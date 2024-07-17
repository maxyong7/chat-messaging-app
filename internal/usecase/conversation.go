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

// MessageRequest struct to hold message data
type MessageRequest struct {
	MessageType string `json:"message_type"`
	Data        Data   `json:"data"`
}

type Data struct {
	MessageUUID      string `json:"message_uuid, omitempty"`
	SenderUUID       string `json:"sender_uuid"`
	ConversationUUID string `json:"conversation_uuid"`
	ReactionData
	MessageData
}
type ReactionData struct {
	ReactionType string `json:"reaction_type"`
}

type MessageData struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageResponse struct {
	MessageType  string       `json:"message_type"`
	ResponseData ResponseData `json:"data"`
}

type ResponseData struct {
	MessageUUID      string `json:"message_uuid, omitempty"`
	ConversationUUID string `json:"conversation_uuid"`
	ResponseReaction
	ResponseMessage
	ResponseError
}

type ResponseReaction struct {
	Reaction string `json:"reaction, omitempty"`
}

type ResponseMessage struct {
	SenderFirstName string    `json:"sender_first_name"`
	SenderLastName  string    `json:"sender_last_name"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
}

type ResponseError struct {
	ErrorMessage string `json:"error_msg"`
	SenderUUID   string `json:"sender_uuid"`
}

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
	Broadcast   chan MessageResponse
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
		Broadcast:   make(chan MessageResponse),
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
			h.HandleMessage(message)
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
func (h *Hub) HandleMessage(message MessageResponse) {
	fmt.Println("HandleMessage: ", message)
	//Check if the message is a type of "message"
	// if message.Type == "message" {
	clients := h.Clients[message.ResponseData.ConversationUUID]
	if message.MessageType == "error" {
		for client := range clients {
			if client.UserInfo.UserUUID == message.ResponseData.SenderUUID {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients[message.ResponseData.ConversationUUID], client)
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
		var msg MessageRequest
		err := c.Conn.ReadJSON(&msg)
		msg.Data.SenderUUID = c.UserInfo.UserUUID
		fmt.Println("readPump: ", msg)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		switch msg.MessageType {
		case sendMessageType:
			msg.Data.ConversationUUID = c.ID
			msg.Data.MessageUUID = uuid.New().String()
			msg.Data.CreatedAt = time.Now()
			err = c.repo.StoreConversation(msg)
			if err != nil {
				fmt.Println("Conversation - readPump - StoreConversation err: ", err)
				errorMsg := c.buildErrorMessage(msg, errProcessingMessage)
				c.hub.Broadcast <- errorMsg
				break
			}
			msgResponse := c.buildMessageResponse(msg)
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

func (c *Client) buildMessageResponse(message MessageRequest) MessageResponse {
	return MessageResponse{
		MessageType: message.MessageType,
		ResponseData: ResponseData{
			MessageUUID:      message.Data.MessageUUID,
			ConversationUUID: message.Data.ConversationUUID,
			ResponseReaction: ResponseReaction{
				Reaction: message.Data.ReactionType,
			},
			ResponseMessage: ResponseMessage{
				SenderFirstName: c.UserInfo.FirstName,
				SenderLastName:  c.UserInfo.LastName,
				Content:         message.Data.Content,
				CreatedAt:       message.Data.CreatedAt,
			},
		},
	}
}

func (c *Client) buildErrorMessage(message MessageRequest, errorMsg string) MessageResponse {
	return MessageResponse{
		MessageType: errorMessageType,
		ResponseData: ResponseData{
			MessageUUID:      message.Data.MessageUUID,
			ConversationUUID: message.Data.ConversationUUID,
			ResponseError: ResponseError{
				ErrorMessage: errorMsg,
				SenderUUID:   c.UserInfo.UserUUID,
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

// func (uc *ConversationUseCase) ServeWsWithRW(w http.ResponseWriter, r *http.Request, hub *Hub) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	// clientID := c.Request.URL.Query().Get("id")
// 	clientID := "thisIsClientID"
// 	// clientID := r.URL.Query().Get("id")

// 	client := NewClient(clientID, conn, hub)
// 	// client := &Client{ID: clientID, Conn: conn, send: make(chan []byte, 256), hub: hub}
// 	// hub.Register <- client
// 	client.hub.Register <- client

// 	go client.writePump()
// 	go client.readPump()
// }
