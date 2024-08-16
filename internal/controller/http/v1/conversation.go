package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/maxyong7/chat-messaging-app/internal/boundary"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type conversationRoutes struct {
	conv     usecase.Conversation
	up       usecase.UserProfile
	msg      usecase.Message
	reaction usecase.Reaction
	l        logger.Interface
}

// Handles api routes for conversation functionality
func newConversationRoute(handler *gin.RouterGroup, c usecase.Conversation, up usecase.UserProfile, msg usecase.Message, reaction usecase.Reaction, l logger.Interface) {
	route := &conversationRoutes{c, up, msg, reaction, l}
	// Initialize a hub and run it with a new thread for websocket connection
	hub := NewHub()
	go hub.Run()

	// Group the routes under the "/conversation" path.
	h := handler.Group("/conversation")
	{
		// Define the endpoints for the conversation functionality.
		h.GET("", route.getConversations)
		h.GET("/ws/:conversationId", route.serveWsController(hub))
	}
}

func (r *conversationRoutes) getConversations(c *gin.Context) {
	// Get decoded 'cursor' value from URL query
	cursor, err := queryParamCursor(c)
	if err != nil {
		// Logs error message
		r.l.Error(err, "http - v1 - getConversations - cursor validation error")
		handleCustomErrors(c, err)
	}

	// Get 'limit' value from URL query and convert into integer type
	// If not value was provided, default to 20
	limit := queryParamInt(c, "limit", 20)

	// Get user_uuid from context
	userId, err := getUserUUIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Build request params entity object
	requestParams := entity.RequestParams{
		Cursor: cursor,
		Limit:  limit,
		UserID: userId,
	}

	// Calls GetConversationList method from conversation entity object
	conversations, err := r.conv.GetConversationList(c.Request.Context(), requestParams)
	if err != nil {
		// Logs the error
		r.l.Error(err, "http - v1 - getConversations - getConversations")

		// If its a known defined error, it writes the status code and return a JSON body with error field accordingly.
		// Else, it defaults to 500 status code and returns 'internal server error' in error field of the JSON body
		handleCustomErrors(c, err)
		return
	}

	var encodedCursor string
	if len(conversations) == limit {
		// Use encoded timestamp from the last message as the cursor for pagination.
		// Boundary object will later provide this value to get the next paginated list
		encodedCursor = encodeCursor(conversations[len(conversations)-1].LastMessageCreatedAt)
	}

	// Build boundary object
	conversationResp := boundary.GetConversationsResponseModel{
		Data: boundary.GetConversationsData{
			Conversations: conversations,
		},
		Pagination: boundary.Pagination{
			Cursor: encodedCursor,
			Limit:  limit,
		},
	}

	// Writes the status code provided in the argument.
	// It also writes a JSON body using the boundary object.
	c.JSON(http.StatusOK, conversationResp)
}

type Client struct {
	ID       string
	UserInfo entity.UserProfile
	Conn     *websocket.Conn
	send     chan boundary.ConversationResponseModel
	hub      *Hub
	route    *conversationRoutes
}

// Method that act as a websocket controller
func (r *conversationRoutes) serveWsController(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user_uuid from context
		userUUID, err := getUserUUIDFromContext(c)
		if err != nil {
			errorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		// Upgrades the HTTP server connection to the WebSocket protocol.
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			// Logs the error
			log.Println(err)
			return
		}

		// Calls GetUserProfile method from user profile entity object
		userInfo, err := r.up.GetUserProfile(c.Request.Context(), userUUID)
		if err != nil || userInfo.UserUUID == "" {
			// Logs the error
			log.Println("ServeWs - GetUserInfo err:", err)
			return
		}

		// Get conversationId from URL parameter
		clientId := c.Param("conversationId")
		// Creates a new client
		client := NewClient(clientId, userInfo, conn, hub, r)
		// Register the new client to the hub.
		// It checks if room exists based on the conversationId, create it if doesn't exist and add client to it
		hub.Register <- client

		// Create a new thread to write and send messages
		go client.writePump()
		// Create a new thread to read any incoming messages
		go client.readPump()
	}
}

type Hub struct {
	Clients     map[string]map[*Client]bool
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan boundary.ConversationResponseModel
	HandleError chan boundary.ConversationResponseModel
	mu          sync.Mutex
}

// Method to initialize a new hub
func NewHub() *Hub {
	return &Hub{
		Clients:     make(map[string]map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan boundary.ConversationResponseModel),
		HandleError: make(chan boundary.ConversationResponseModel),
	}
}

// Method to run the hub
func (h *Hub) Run() {
	// Keeps running until the application stops
	for {
		select {
		// Register a client if 'Register' is called
		case client := <-h.Register:
			// Locks mutex
			h.mu.Lock()

			// Register the new client to the hub.
			// It checks if room exists based on the conversationId, create it if doesn't exist and add client to it
			h.RegisterNewClient(client)

			// Unlocks mutex
			h.mu.Unlock()

			// Logs when a client has connected to the hub's room
			fmt.Printf("Client %s connected\n", client.ID)

		// Unregister a client if 'Unregister' is called
		case client := <-h.Unregister:
			// Locks mutex
			h.mu.Lock()

			// Find the client based on client's ID
			if _, ok := h.Clients[client.ID]; ok {
				// If found delete the client id key stored in the hub
				delete(h.Clients, client.ID)

				// Closes the client websocket channel
				close(client.send)
			}

			// Unlocks mutex
			h.mu.Unlock()

			// Logs when a client has disconnected from the hub
			fmt.Printf("Client %s disconnected\n", client.ID)

		// Broadcast messages to client(s) if 'Broadcast' is called
		case message := <-h.Broadcast:
			// Locks mutex
			h.mu.Lock()

			// Method to handle broadcasting message based on type of message
			h.HandleBroadcast(message)

			// Unlocks mutex
			h.mu.Unlock()
		}
	}
}

// This method checks if room exists and if not create it and add client to it
func (h *Hub) RegisterNewClient(client *Client) {
	// Check if connection already exist in the hub's dictionary based on client id
	connections := h.Clients[client.ID]
	if connections == nil {
		// If connection does not exist, create a connection based on client info
		connections = make(map[*Client]bool)
		// Add connection into the hub's dictionary
		h.Clients[client.ID] = connections
	}
	// Set hub dictionary value to true
	h.Clients[client.ID][client] = true

	// Logs the size of client stored in hub's dictionary
	fmt.Println("Size of Clients: ", len(h.Clients[client.ID]))
}

// Method to handle broadcasting message based on type of message
func (h *Hub) HandleBroadcast(message boundary.ConversationResponseModel) {
	// Logs message that being processed
	fmt.Println("HandleBroadcast: ", message)

	// Get all the clients connected to the same ConversationUUID
	clients := h.Clients[message.Data.ConversationUUID]

	// If message type is 'error', broadcast the error message only to the sender
	if message.MessageType == "error" {
		// Loops through the list of clients to find the sender
		for client := range clients {
			if client.UserInfo.UserUUID == message.Data.SenderUUID {
				select {
				// Sends the error message to the sender via websocket
				case client.send <- message:

				// For any unexpected case, close websocket and delete client from hub
				default:
					close(client.send)
					delete(h.Clients[message.Data.ConversationUUID], client)
				}
			}
		}
		return
	}

	// Loops through the list of clients
	for client := range clients {
		select {
		//Send message to each client according to the ConversationUUID
		case client.send <- message:

		// For any unexpected case, close websocket and delete client from hub
		default:
			close(client.send)
			delete(h.Clients[message.Data.ConversationUUID], client)
		}
	}

}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	errorMessageType          = "error"
	sendMessageType           = "send_message"
	addReactionMessageType    = "add_reaction"
	removeReactionMessageType = "remove_reaction"
	deleteMessageType         = "delete_message"
	errProcessingMessage      = "error processing message"
	errProcessingReaction     = "error processing reaction"
	errOnlyAuthorCanDeleteMsg = "cannot delete because user is not message author"
)

// readPump handles reading messages from the WebSocket connection.
func (c *Client) readPump() {
	// Ensure the connection is closed and the client is unregistered from the hub when the function exits.
	defer func() {
		c.hub.Unregister <- c
		c.Conn.Close()
	}()

	// Set the maximum size for incoming messages.
	c.Conn.SetReadLimit(maxMessageSize)

	// Set the initial read deadline for incoming messages, which is updated upon receiving Pong messages.
	// If a read has timed out, the websocket connection state is corrupt and all future reads will return an error
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	// Define a handler for Pong messages to extend the read deadline.
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start reading messages in a loop.
	for {
		var msgReq boundary.ConversationRequestModel

		// Read a JSON message from the WebSocket connection into msgReq.
		err := c.Conn.ReadJSON(&msgReq)
		fmt.Println("readPump: ", msgReq)

		// If there is an error reading the message, log the error and break the loop.
		if err != nil {
			fmt.Println("readPump - Error ReadJSON", err)
			break
		}

		// Handle the conversation message using the client's handleConversation method.
		c.handleConversation(msgReq, c.UserInfo)
	}
}

// writePump handles writing messages to the WebSocket connection.
func (c *Client) writePump() {
	// Create a new ticker that triggers periodically based on pingPeriod.
	ticker := time.NewTicker(pingPeriod)

	// Ensure the ticker is stopped and the connection is closed when the function exits.
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	// Start a loop to handle sending messages and pings.
	for {
		select {
		case message, ok := <-c.send:
			// Set a write deadline for sending messages.
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// If the channel was closed by the hub, send a WebSocket close message and return.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				// Write the JSON message to the WebSocket connection.
				err := c.Conn.WriteJSON(message)
				fmt.Println("writePump: ", message)

				// If there is an error writing the message, log the error and break the loop.
				if err != nil {
					fmt.Println("Error: ", err)
					break
				}
			}
		case <-ticker.C:
			// When the ticker triggers, send a Ping message to keep the connection alive.
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Method that creates a new client
func NewClient(id string, userInfo entity.UserProfile, conn *websocket.Conn, hub *Hub, route *conversationRoutes) *Client {
	return &Client{
		ID:       id,
		UserInfo: userInfo,
		Conn:     conn,
		send:     make(chan boundary.ConversationResponseModel, 256),
		hub:      hub,
		route:    route,
	}
}

// Initialize a websocket upgrader as a variable
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handleConversation processes different types of conversation requests such as sending, deleting,
// adding a reaction to a message, or removing a reaction. It also handles errors and broadcasts
// appropriate messages to the hub.
func (c *Client) handleConversation(convReq boundary.ConversationRequestModel, userInfo entity.UserProfile) {
	// Extract sender and conversation identifiers from userInfo and the client instance.
	senderUUID := userInfo.UserUUID
	conversationUUID := c.ID

	// Create a background context
	ctx := context.Background()

	// Handle different types of conversation requests based on the MessageType in convReq.
	switch convReq.MessageType {
	case sendMessageType:
		// Unmarshal the data in convReq into a ChatInterface boundary object.
		var sendMessageRequest boundary.ChatInterface
		err := json.Unmarshal(convReq.Data, &sendMessageRequest)
		if err != nil {
			// If there's an error in unmarshalling, log it and broadcast an error message.
			fmt.Println("handleConversation - unmarshall error for sendMessageRequest", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}

		// Create a new conversation entity with the provided data.
		conv := entity.Conversation{
			SenderUUID:       userInfo.UserUUID,
			ConversationUUID: c.ID,
			MessageUUID:      uuid.New().String(),
			Content:          sendMessageRequest.Content,
			CreatedAt:        time.Now(),
		}

		// Store the conversation and message by calling conversation entity object's StoreConversationAndMessage method
		err = c.route.conv.StoreConversationAndMessage(ctx, conv)
		if err != nil {
			// If there's an error storing the message, log it and broadcast an error message.
			fmt.Println("Conversation - handleConversation - StoreConversation err: ", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}

		// Build a response message and broadcast it.
		sendMsgResponse := buildSendMessageResponse(conv, userInfo)
		c.hub.Broadcast <- sendMsgResponse

	case deleteMessageType:
		// Unmarshal the data in convReq into a DeleteMessageRequest object.
		var deleteMessageRequest boundary.DeleteMessageRequest
		err := json.Unmarshal(convReq.Data, &deleteMessageRequest)
		if err != nil {
			// If there's an error in unmarshalling, log it and broadcast an error message.
			fmt.Println("handleConversation - unmarshall error for deleteMessageRequest", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}

		// Create a message entity for deletion based on the request data.
		msg := entity.Message{
			SenderUUID:  senderUUID,
			MessageUUID: deleteMessageRequest.MessageUUID,
		}

		// Attempt to delete the message by calling message entity object's DeleteMessage method.
		valid, err := c.route.msg.DeleteMessage(ctx, msg)
		if err != nil {
			fmt.Println("Conversation - readPump - DeleteMessage err: ", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}

		// If deletion was not valid (e.g., the user is not the author), broadcast an error message.
		if !valid {
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errOnlyAuthorCanDeleteMsg)
			c.hub.Broadcast <- errorMsg
		}

		// Build a response for the deletion and broadcast it.
		deleteMsgResponse := buildDeleteMessageResponse(msg, conversationUUID)
		c.hub.Broadcast <- deleteMsgResponse

	case addReactionMessageType:
		// Unmarshal the data in convReq into a MessageReactionMenu object.
		var addReactionRequest boundary.MessageReactionMenu
		err := json.Unmarshal(convReq.Data, &addReactionRequest)
		if err != nil {
			// If there's an error in unmarshalling, log it and broadcast an error message.
			fmt.Println("handleConversation - unmarshall error for addReactionRequest", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}

		// Create a reaction entity with the provided data.
		reaction := entity.Reaction{
			MessageUUID:  addReactionRequest.MessageUUID,
			SenderUUID:   userInfo.UserUUID,
			ReactionType: addReactionRequest.ReactionType,
		}

		// Store the reaction by calling reaction entity object's StoreReaction method.
		err = c.route.reaction.StoreReaction(ctx, reaction)
		if err != nil {
			// If there's an error storing the reaction, log it and broadcast an error message.
			fmt.Println("Conversation - readPump - StoreReaction err: ", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}

		// Build a response for adding the reaction and broadcast it.
		addReactionResponse := buildReactionResponse(addReactionMessageType, reaction, conversationUUID)
		c.hub.Broadcast <- addReactionResponse

	case removeReactionMessageType:
		// Unmarshal the data in convReq into a RemoveReactionRequest object.
		var removeReactionRequest boundary.RemoveReactionRequest
		err := json.Unmarshal(convReq.Data, &removeReactionRequest)
		if err != nil {
			// If there's an error in unmarshalling, log it and broadcast an error message.
			fmt.Println("handleConversation - unmarshall error for removeReactionRequest", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}

		// Create a reaction entity for removal based on the request data.
		reaction := entity.Reaction{
			MessageUUID: removeReactionRequest.MessageUUID,
			SenderUUID:  userInfo.UserUUID,
		}

		// Remove the reaction by calling reaction entity object's RemoveReaction method.
		err = c.route.reaction.RemoveReaction(ctx, reaction)
		if err != nil {
			// If there's an error removing the reaction, log it and broadcast an error message.
			fmt.Println("Conversation - readPump - RemoveReaction err: ", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}

		// Build a response for removing the reaction and broadcast it.
		removeReactionResponse := buildReactionResponse(removeReactionMessageType, reaction, conversationUUID)
		c.hub.Broadcast <- removeReactionResponse
	}
}

// Method to build send message response body
func buildSendMessageResponse(conv entity.Conversation, userInfo entity.UserProfile) boundary.ConversationResponseModel {
	return boundary.ConversationResponseModel{
		MessageType: sendMessageType,
		Data: boundary.ConversationResponseData{
			SenderUUID:       conv.SenderUUID,
			ConversationUUID: conv.ConversationUUID,
			SendMessageResponseData: boundary.SendMessageResponseData{
				SenderFirstName: userInfo.FirstName,
				SenderLastName:  userInfo.LastName,
				SenderAvatar:    userInfo.Avatar,
				Content:         conv.Content,
				MessageUUID:     conv.MessageUUID,
				CreatedAt:       conv.CreatedAt,
			},
		},
	}
}

// Method to build delete message response body
func buildDeleteMessageResponse(msg entity.Message, conversationUUID string) boundary.ConversationResponseModel {
	return boundary.ConversationResponseModel{
		MessageType: deleteMessageType,
		Data: boundary.ConversationResponseData{
			SenderUUID:       msg.SenderUUID,
			ConversationUUID: conversationUUID,
			MessageDeletionConfirmation: boundary.MessageDeletionConfirmation{
				MessageUUID: msg.MessageUUID,
			},
		},
	}
}

// Method to build reaction response body
func buildReactionResponse(reactionType string, reaction entity.Reaction, conversationUUID string) boundary.ConversationResponseModel {
	return boundary.ConversationResponseModel{
		MessageType: reactionType,
		Data: boundary.ConversationResponseData{
			SenderUUID:       reaction.SenderUUID,
			ConversationUUID: conversationUUID,
			ReactionResponseData: boundary.ReactionResponseData{
				MessageUUID: reaction.MessageUUID,
				Reaction:    reaction.ReactionType,
			},
		},
	}
}

// Method to build error message response body
func (c *Client) buildErrorMessage(senderUUID string, conversationUUID string, errorMsg string) boundary.ConversationResponseModel {
	return boundary.ConversationResponseModel{
		MessageType: errorMessageType,
		Data: boundary.ConversationResponseData{
			SenderUUID:       senderUUID,
			ConversationUUID: conversationUUID,
			ErrorResponseData: boundary.ErrorResponseData{
				ErrorMessage: errorMsg,
			},
		},
	}
}
