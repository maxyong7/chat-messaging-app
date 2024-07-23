package usecase

import (
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
)

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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// ConversationUseCase -.
type ConversationUseCase struct {
	repo         ConversationRepo
	userRepo     UserRepo
	reactionRepo ReactionRepo
	messageRepo  MessageRepo
	// webAPI  TranslationWebAPI
}

// // MessageRequest struct to hold message data
// type MessageRequest struct {
// 	MessageType string `json:"message_type"`
// 	Data        Data   `json:"data"`
// }

//	type Data struct {
//		MessageUUID      string `json:"message_uuid,omitempty"`
//		SenderUUID       string `json:"sender_uuid"`
//		ConversationUUID string `json:"conversation_uuid"`
//		ReactionData
//		MessageData
//	}
type ReactionData struct {
	ReactionType string `json:"reaction_type"`
}

// type MessageData struct {
// 	Content   string    `json:"content"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// type MessageResponse struct {
// 	MessageType  string       `json:"message_type"`
// 	ResponseData ResponseData `json:"data"`
// }

// type ResponseData struct {
// 	MessageUUID      string `json:"message_uuid,omitempty"`
// 	ConversationUUID string `json:"conversation_uuid"`
// 	ResponseReaction
// 	ResponseMessage
// 	ResponseError
// }

// type ResponseReaction struct {
// 	Reaction string `json:"reaction,omitempty"`
// }

// type ResponseMessage struct {
// 	SenderFirstName string    `json:"sender_first_name"`
// 	SenderLastName  string    `json:"sender_last_name"`
// 	Content         string    `json:"content"`
// 	CreatedAt       time.Time `json:"created_at"`
// }

// type ResponseError struct {
// 	ErrorMessage string `json:"error_msg"`
// 	SenderUUID   string `json:"sender_uuid"`
// }

type Client struct {
	ID           string
	UserInfo     *entity.UserInfoDTO
	Conn         *websocket.Conn
	send         chan MessageResponse
	hub          *Hub
	repo         ConversationRepo
	reactionRepo ReactionRepo
	messageRepo  MessageRepo
}

type Hub struct {
	Clients     map[string]map[*Client]bool
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan boundary.ConversationResponseModel
	HandleError chan MessageResponse
	mu          sync.Mutex
}

// New -.
func NewConversation(r ConversationRepo, userRepo UserRepo, reactionRepo ReactionRepo, msgRepo MessageRepo) *ConversationUseCase {
	return &ConversationUseCase{
		repo:         r,
		userRepo:     userRepo,
		reactionRepo: reactionRepo,
		messageRepo:  msgRepo,
	}
}

func NewHub() *Hub {
	return &Hub{
		Clients:     make(map[string]map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan boundary.ConversationResponseModel),
		HandleError: make(chan MessageResponse),
	}
}

// var hub = Hub{
// 	Clients:    make(map[string]*Client),
// 	Register:   make(chan *Client),
// 	Unregister: make(chan *Client),
// 	Broadcast:  make(chan []byte),
// }

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

// function to remove client from room
func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.Clients[client.ID]; ok {
		delete(h.Clients[client.ID], client)
		close(client.send)
		fmt.Println("Removed client")
	}
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
			delete(h.Clients[message.ResponseData.ConversationUUID], client)
		}
	}

}

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
		c.handleConversation(msgReq, c.UserInfo.UserUUID)
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

// NewClient creates a new client
func NewClient(id string, userInfo *entity.UserInfoDTO, conn *websocket.Conn, hub *Hub, uc ConversationUseCase) *Client {
	return &Client{
		ID:           id,
		UserInfo:     userInfo,
		Conn:         conn,
		send:         make(chan MessageResponse, 256),
		hub:          hub,
		repo:         uc.repo,
		reactionRepo: uc.reactionRepo,
		messageRepo:  uc.messageRepo,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (uc *ConversationUseCase) ServeWs(c *gin.Context, hub *Hub, userUuid string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// clientId := c.Request.URL.Query().Get("clientId")
	// clientID := "thisIsClientID"
	// clientID := c.Request.Header.Get("Sec-Websocket-Key")
	userInfo, err := uc.userRepo.GetUserInfo(c.Request.Context(), userUuid)
	if err != nil || userInfo == nil {
		log.Println("ServeWs - GetUserInfo err:", err)
		return
	}

	clientId := c.Param("conversationId")
	client := NewClient(clientId, userInfo, conn, hub, *uc)
	hub.Register <- client
	// client.Hub.Register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) handleConversation(convReq boundary.ConversationRequestModel, senderUUID string) {
	conv := entity.Conversation{
		SenderUUID:       senderUUID,
		ConversationUUID: c.ID,
	}
	switch convReq.MessageType {
	case sendMessageType:
		msg := entity.Message{
			MessageUUID: uuid.New().String(),
			Content:     convReq.Data.SendMessageRequest.Content,
			CreatedAt:   time.Now(),
		}
		err := c.repo.StoreConversation(conv, msg)
		if err != nil {
			fmt.Println("Conversation - handleConversation - StoreConversation err: ", err)
			errorMsg := c.buildErrorMessage(conv, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}
		msgResponse := conversationMessageToResponse(msg, sendMessageType)
		c.hub.Broadcast <- msgResponse
	case deleteMessageType:
		msg.Data.ConversationUUID = c.ID
		valid, err := c.messageRepo.ValidateMessageSentByUser(msg)
		if err != nil {
			fmt.Println("Conversation - readPump - ValidateMessageSentByUser err: ", err)
			errorMsg := c.buildErrorMessage(msg, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}
		if !valid {
			errorMsg := c.buildErrorMessage(msg, errOnlyAuthorCanDeleteMsg)
			c.hub.Broadcast <- errorMsg
		}
		err = c.messageRepo.DeleteMessage(msg)
		if err != nil {
			fmt.Println("Conversation - readPump - DeleteMessage err: ", err)
			errorMsg := c.buildErrorMessage(msg, errProcessingMessage)
			c.hub.Broadcast <- errorMsg
			break
		}
		msgResponse := c.buildMessageResponse(msg)
		c.hub.Broadcast <- msgResponse
	case addReactionMessageType:
		err = c.reactionRepo.StoreReaction(msg)
		if err != nil {
			fmt.Println("Conversation - readPump - StoreReaction err: ", err)
			errorMsg := c.buildErrorMessage(msg, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}
		msgResponse := c.buildMessageResponse(msg)
		c.hub.Broadcast <- msgResponse
	case removeReactionMessageType:
		err = c.reactionRepo.RemoveReaction(msg)
		if err != nil {
			fmt.Println("Conversation - readPump - RemoveReaction err: ", err)
			errorMsg := c.buildErrorMessage(msg, errProcessingReaction)
			c.hub.Broadcast <- errorMsg
			break
		}
		msgResponse := c.buildMessageResponse(msg)
		c.hub.Broadcast <- msgResponse
	}
}

func conversationMessageToResponse(convMsg entity.ConversationMessage, messageType string) boundary.ConversationResponseModel {
	return boundary.ConversationResponseModel{
		MessageType: messageType,
		Data: boundary.ConversationResponseData{
			Conversation: entity.Conversation{
				SenderUUID:          "",
				ConversationUUID:    "",
				ConversationMessage: entity.ConversationMessage{},
			},
		},
	}
}
