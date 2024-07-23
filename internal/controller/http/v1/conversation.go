package v1

import (
	"context"
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
		h.GET("/ws/:conversationId", route.ServeWsController(hub))
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

type Client struct {
	ID       string
	UserInfo entity.UserInfo
	Conn     *websocket.Conn
	send     chan boundary.ConversationResponseModel
	hub      *Hub
	route    *conversationRoutes
}

func (r *conversationRoutes) ServeWsController(hub *Hub) gin.HandlerFunc {
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

		userInfo, err := r.up.GetUserInfo(c.Request.Context(), userUUID)
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
func NewClient(id string, userInfo entity.UserInfo, conn *websocket.Conn, hub *Hub, route *conversationRoutes) *Client {
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

func (c *Client) handleConversation(convReq boundary.ConversationRequestModel, userInfo entity.UserInfo) {
	conv := entity.Conversation{
		SenderUUID:       userInfo.UserUUID,
		ConversationUUID: c.ID,
	}
	ctx := context.Background()
	switch convReq.MessageType {
	case sendMessageType:
		msg := entity.Message{
			MessageUUID: uuid.New().String(),
			Content:     convReq.Data.SendMessageRequest.Content,
			CreatedAt:   time.Now(),
		}
		err := c.route.conv.StoreConversation(ctx, conv, msg)
		if err != nil {
			fmt.Println("Conversation - handleConversation - StoreConversation err: ", err)
			errorMsg := c.buildErrorMessage(conv, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}
		sendMsgResponse := buildSendMessageResponse(conv, msg, userInfo)
		c.hub.Broadcast <- sendMsgResponse
	case deleteMessageType:
		msg := entity.Message{
			MessageUUID: convReq.Data.DeleteMessageRequest.MessageUUID,
		}
		valid, err := c.route.msg.ValidateMessageSentByUser(ctx, conv, msg)
		if err != nil {
			fmt.Println("Conversation - readPump - ValidateMessageSentByUser err: ", err)
			errorMsg := c.buildErrorMessage(conv, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}
		if !valid {
			errorMsg := c.buildErrorMessage(conv, errOnlyAuthorCanDeleteMsg)
			c.hub.Broadcast <- errorMsg
		}
		err = c.route.msg.DeleteMessage(ctx, conv, msg)
		if err != nil {
			fmt.Println("Conversation - readPump - DeleteMessage err: ", err)
			errorMsg := c.buildErrorMessage(conv, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}
		deleteMsgResponse := buildDeleteMessageResponse(conv, msg)
		c.hub.Broadcast <- deleteMsgResponse
	case addReactionMessageType:
		reaction := entity.Reaction{
			MessageUUID:  convReq.Data.AddReactionRequest.MessageUUID,
			SenderUUID:   userInfo.UserUUID,
			ReactionType: convReq.Data.AddReactionRequest.ReactionType,
		}
		err := c.route.reaction.StoreReaction(ctx, reaction)
		if err != nil {
			fmt.Println("Conversation - readPump - StoreReaction err: ", err)
			errorMsg := c.buildErrorMessage(conv, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}
		addReactionResponse := buildReactionResponse(addReactionMessageType, conv, reaction)
		c.hub.Broadcast <- addReactionResponse
	case removeReactionMessageType:
		reaction := entity.Reaction{
			MessageUUID:  convReq.Data.AddReactionRequest.MessageUUID,
			SenderUUID:   userInfo.UserUUID,
			ReactionType: convReq.Data.AddReactionRequest.ReactionType,
		}
		err := c.route.reaction.RemoveReaction(ctx, reaction)
		if err != nil {
			fmt.Println("Conversation - readPump - RemoveReaction err: ", err)
			errorMsg := c.buildErrorMessage(conv, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}
		removeReactionResponse := buildReactionResponse(removeReactionMessageType, conv, reaction)
		c.hub.Broadcast <- removeReactionResponse
	}
}

func buildSendMessageResponse(conv entity.Conversation, msg entity.Message, userInfo entity.UserInfo) boundary.ConversationResponseModel {
	return boundary.ConversationResponseModel{
		MessageType: sendMessageType,
		Data: boundary.ConversationResponseData{
			SenderUUID:       conv.SenderUUID,
			ConversationUUID: conv.ConversationUUID,
			SendMessageResponseData: boundary.SendMessageResponseData{
				SenderFirstName: userInfo.FirstName,
				SenderLastName:  userInfo.LastName,
				Content:         msg.Content,
				MessageUUID:     msg.MessageUUID,
				CreatedAt:       msg.CreatedAt,
			},
		},
	}
}

func buildDeleteMessageResponse(conv entity.Conversation, msg entity.Message) boundary.ConversationResponseModel {
	return boundary.ConversationResponseModel{
		MessageType: deleteMessageType,
		Data: boundary.ConversationResponseData{
			SenderUUID:       conv.SenderUUID,
			ConversationUUID: conv.ConversationUUID,
			DeleteMessageResponseData: boundary.DeleteMessageResponseData{
				MessageUUID: msg.MessageUUID,
			},
		},
	}
}

func buildReactionResponse(reactionType string, conv entity.Conversation, reaction entity.Reaction) boundary.ConversationResponseModel {
	return boundary.ConversationResponseModel{
		MessageType: reactionType,
		Data: boundary.ConversationResponseData{
			SenderUUID:       conv.SenderUUID,
			ConversationUUID: conv.ConversationUUID,
			AddReactionResponseData: boundary.ReactionResponseData{
				MessageUUID: reaction.MessageUUID,
				Reaction:    reaction.ReactionType,
			},
		},
	}
}

func (c *Client) buildErrorMessage(conv entity.Conversation, errorMsg string) boundary.ConversationResponseModel {
	return boundary.ConversationResponseModel{
		MessageType: errorMessageType,
		Data: boundary.ConversationResponseData{
			SenderUUID:       conv.SenderUUID,
			ConversationUUID: conv.ConversationUUID,
			ErrorResponseData: boundary.ErrorResponseData{
				ErrorMessage: errorMsg,
			},
		},
	}
}
