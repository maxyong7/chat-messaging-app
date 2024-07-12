package usecase

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// ConversationUseCase -.
type ConversationUseCase struct {
	repo   ConversationRepo
	webAPI TranslationWebAPI
}

// Message struct to hold message data
type Message struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	ID        string `json:"id"`
}

type Client struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	send   chan Message
	hub    *Hub
}

type Hub struct {
	Clients    map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	mu         sync.Mutex
}

// New -.
func NewConversation(r ConversationRepo, w TranslationWebAPI) *ConversationUseCase {
	return &ConversationUseCase{
		repo:   r,
		webAPI: w,
	}
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]map[*Client]bool),
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
func (h *Hub) HandleMessage(message Message) {
	fmt.Println("HandleMessage: ", message)
	//Check if the message is a type of "message"
	if message.Type == "message" {
		clients := h.Clients[message.ID]
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.Clients[message.ID], client)
			}
		}
	}

	//Check if the message is a type of "notification"
	if message.Type == "notification" {
		fmt.Println("Notification: ", message.Content)
		clients := h.Clients[message.Recipient]
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.Clients[message.Recipient], client)
			}
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
		var msg Message
		err := c.Conn.ReadJSON(&msg)

		fmt.Println("readPump: ", msg)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		c.hub.Broadcast <- msg
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
func NewClient(id string, userId string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{ID: id, UserID: userId, Conn: conn, send: make(chan Message, 256), hub: hub}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (uc *ConversationUseCase) ServeWs(c *gin.Context, hub *Hub, userId string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// clientId := c.Request.URL.Query().Get("clientId")
	// clientID := "thisIsClientID"
	// clientID := c.Request.Header.Get("Sec-Websocket-Key")

	clientId := c.Param("conversationId")
	client := NewClient(clientId, userId, conn, hub)
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
