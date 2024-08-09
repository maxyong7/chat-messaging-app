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

func newConversationRoute(handler *gin.RouterGroup, c usecase.Conversation, up usecase.UserProfile, msg usecase.Message, reaction usecase.Reaction, l logger.Interface) {
	route := &conversationRoutes{c, up, msg, reaction, l}
	hub := NewHub()
	go hub.Run()

	h := handler.Group("/conversation")
	{
		h.GET("", route.getConversations)
		h.GET("/ws/:conversationId", route.serveWsController(hub))
		// h.GET("/ws/:clientId", func(c *gin.Context) {

		// 	route.t.ServeWs(c, hub)
		// })
		// http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		// 	route.t.ServeWsWithRW(w, r, hub)
		// })
		// http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		// 	route.t.ServeWsWithRW(w, r, hub)
		// })
		// h.GET("/ws", route.ServeWsController)
	}
}

func (r *conversationRoutes) getConversations(c *gin.Context) {
	cursor, err := queryParamCursor(c)
	if err != nil {
		r.l.Error(err, "http - v1 - getConversations - cursor validation error")
		handleCustomErrors(c, err)
	}
	limit := queryParamInt(c, "limit", 20)

	userId, err := getUserUUIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	requestParams := entity.RequestParams{
		Cursor: cursor,
		Limit:  limit,
		UserID: userId,
	}

	conversations, err := r.conv.GetConversationList(c.Request.Context(), requestParams)
	if err != nil {
		r.l.Error(err, "http - v1 - getConversations - getConversations")
		handleCustomErrors(c, err)
		return
	}

	var encodedCursor string
	if len(conversations) == limit {
		encodedCursor = encodeCursor(conversations[len(conversations)-1].LastMessageCreatedAt)
	}

	conversationResp := boundary.GetConversationsResponseModel{
		Data: boundary.GetConversationsData{
			Conversations: conversations,
		},
		Pagination: boundary.Pagination{
			Cursor: encodedCursor,
			Limit:  limit,
		},
	}

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

func (r *conversationRoutes) serveWsController(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, err := getUserUUIDFromContext(c)
		if err != nil {
			errorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		userInfo, err := r.up.GetUserProfile(c.Request.Context(), userUUID)
		if err != nil || userInfo.UserUUID == "" {
			log.Println("ServeWs - GetUserInfo err:", err)
			return
		}

		clientId := c.Param("conversationId")
		client := NewClient(clientId, userInfo, conn, hub, r)
		hub.Register <- client

		go client.writePump()
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

func NewHub() *Hub {
	return &Hub{
		Clients:     make(map[string]map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan boundary.ConversationResponseModel),
		HandleError: make(chan boundary.ConversationResponseModel),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.RegisterNewClient(client)
			// h.Clients[client.ID] = client
			h.mu.Unlock()
			fmt.Printf("Client %s connected\n", client.ID)
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.send)
			}
			h.mu.Unlock()
			fmt.Printf("Client %s disconnected\n", client.ID)
		case message := <-h.Broadcast:
			h.mu.Lock()
			h.HandleBroadcast(message)
			// for id, client := range h.Clients {
			// 	select {
			// 	case client.Send <- message:
			// 	default:
			// 		close(client.Send)
			// 		delete(h.Clients, id)
			// 	}
			// }
			h.mu.Unlock()
		}
	}
}

// function check if room exists and if not create it and add client to it
func (h *Hub) RegisterNewClient(client *Client) {
	connections := h.Clients[client.ID]
	if connections == nil {
		connections = make(map[*Client]bool)
		h.Clients[client.ID] = connections
	}
	h.Clients[client.ID][client] = true

	fmt.Println("Size of Clients: ", len(h.Clients[client.ID]))
}

// function to handle message based on type of message
func (h *Hub) HandleBroadcast(message boundary.ConversationResponseModel) {
	fmt.Println("HandleBroadcast: ", message)
	//Check if the message is a type of "message"
	// if message.Type == "message" {
	clients := h.Clients[message.Data.ConversationUUID]
	if message.MessageType == "error" {
		for client := range clients {
			if client.UserInfo.UserUUID == message.Data.SenderUUID {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients[message.Data.ConversationUUID], client)
				}
			}
		}
		return
	}
	for client := range clients {
		select {
		case client.send <- message:
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

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msgReq boundary.ConversationRequestModel
		err := c.Conn.ReadJSON(&msgReq)
		fmt.Println("readPump: ", msgReq)
		if err != nil {
			fmt.Println("readPump - Error ReadJSON", err)
			break
		}
		c.handleConversation(msgReq, c.UserInfo)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {

				err := c.Conn.WriteJSON(message)

				fmt.Println("writePump: ", message)
				if err != nil {
					fmt.Println("Error: ", err)
					break
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// NewClient creates a new client
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Client) handleConversation(convReq boundary.ConversationRequestModel, userInfo entity.UserProfile) {
	senderUUID := userInfo.UserUUID
	conversationUUID := c.ID
	ctx := context.Background()
	switch convReq.MessageType {
	case sendMessageType:
		var sendMessageRequest boundary.SendMessageRequest
		err := json.Unmarshal(convReq.Data, &sendMessageRequest)
		if err != nil {
			fmt.Println("handleConversation - unmarshall error for sendMessageRequest", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}
		conv := entity.Conversation{
			SenderUUID:       userInfo.UserUUID,
			ConversationUUID: c.ID,
			MessageUUID:      uuid.New().String(),
			Content:          sendMessageRequest.Content,
			CreatedAt:        time.Now(),
		}
		err = c.route.conv.StoreConversationAndMessage(ctx, conv)
		if err != nil {
			fmt.Println("Conversation - handleConversation - StoreConversation err: ", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}
		sendMsgResponse := buildSendMessageResponse(conv, userInfo)
		c.hub.Broadcast <- sendMsgResponse
	case deleteMessageType:
		var deleteMessageRequest boundary.DeleteMessageRequest
		err := json.Unmarshal(convReq.Data, &deleteMessageRequest)
		if err != nil {
			fmt.Println("handleConversation - unmarshall error for deleteMessageRequest", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}
		msg := entity.Message{
			SenderUUID:  senderUUID,
			MessageUUID: deleteMessageRequest.MessageUUID,
		}
		valid, err := c.route.msg.DeleteMessage(ctx, msg)
		if err != nil {
			fmt.Println("Conversation - readPump - DeleteMessage err: ", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}
		if !valid {
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errOnlyAuthorCanDeleteMsg)
			c.hub.Broadcast <- errorMsg
		}
		deleteMsgResponse := buildDeleteMessageResponse(msg, conversationUUID)
		c.hub.Broadcast <- deleteMsgResponse
	case addReactionMessageType:
		var addReactionRequest boundary.AddReactionRequest
		err := json.Unmarshal(convReq.Data, &addReactionRequest)
		if err != nil {
			fmt.Println("handleConversation - unmarshall error for addReactionRequest", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}
		reaction := entity.Reaction{
			MessageUUID:  addReactionRequest.MessageUUID,
			SenderUUID:   userInfo.UserUUID,
			ReactionType: addReactionRequest.ReactionType,
		}
		err = c.route.reaction.StoreReaction(ctx, reaction)
		if err != nil {
			fmt.Println("Conversation - readPump - StoreReaction err: ", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}
		addReactionResponse := buildReactionResponse(addReactionMessageType, reaction, conversationUUID)
		c.hub.Broadcast <- addReactionResponse
	case removeReactionMessageType:
		var removeReactionRequest boundary.RemoveReactionRequest
		err := json.Unmarshal(convReq.Data, &removeReactionRequest)
		if err != nil {
			fmt.Println("handleConversation - unmarshall error for removeReactionRequest", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}
		reaction := entity.Reaction{
			MessageUUID: removeReactionRequest.MessageUUID,
			SenderUUID:  userInfo.UserUUID,
		}
		err = c.route.reaction.RemoveReaction(ctx, reaction)
		if err != nil {
			fmt.Println("Conversation - readPump - RemoveReaction err: ", err)
			errorMsg := c.buildErrorMessage(senderUUID, conversationUUID, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}
		removeReactionResponse := buildReactionResponse(removeReactionMessageType, reaction, conversationUUID)
		c.hub.Broadcast <- removeReactionResponse
	}
}

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

func buildDeleteMessageResponse(msg entity.Message, conversationUUID string) boundary.ConversationResponseModel {
	return boundary.ConversationResponseModel{
		MessageType: deleteMessageType,
		Data: boundary.ConversationResponseData{
			SenderUUID:       msg.SenderUUID,
			ConversationUUID: conversationUUID,
			DeleteMessageResponseData: boundary.DeleteMessageResponseData{
				MessageUUID: msg.MessageUUID,
			},
		},
	}
}

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
