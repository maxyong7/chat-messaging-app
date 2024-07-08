package usecase

import (
	"bytes"
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

// New -.
func NewConversation(r ConversationRepo, w TranslationWebAPI) *ConversationUseCase {
	return &ConversationUseCase{
		repo:   r,
		webAPI: w,
	}
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
	Hub  *Hub
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
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
			h.Clients[client.ID] = client
			h.mu.Unlock()
			fmt.Printf("Client %s connected\n", client.ID)
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
			fmt.Printf("Client %s disconnected\n", client.ID)
		case message := <-h.Broadcast:
			h.mu.Lock()
			for id, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, id)
				}
			}
			h.mu.Unlock()
		}
	}
}

// func (c *Client) ReadPump() {
// 	defer func() {
// 		c.Hub.Unregister <- c
// 		c.Conn.Close()
// 	}()
// 	for {
// 		_, message, err := c.Conn.ReadMessage()
// 		if err != nil {
// 			log.Println("read:", err)
// 			break
// 		}
// 		c.Hub.Broadcast <- message
// 	}
// }

// func (c *Client) WritePump() {
// 	for message := range c.Send {
// 		err := c.Conn.WriteMessage(websocket.TextMessage, message)
// 		if err != nil {
// 			log.Println("write:", err)
// 			break
// 		}
// 	}
// 	c.Conn.Close()
// }

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.Hub.Broadcast <- message
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
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (uc *ConversationUseCase) ServeWs(c *gin.Context, hub *Hub) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// clientID := c.Request.URL.Query().Get("id")
	clientID := "thisIsClientID"
	// clientID := r.URL.Query().Get("id")
	client := &Client{ID: clientID, Conn: conn, Send: make(chan []byte, 256), Hub: hub}
	hub.Register <- client
	// client.Hub.Register <- client

	go client.writePump()
	go client.readPump()
}

func (uc *ConversationUseCase) ServeWsWithRW(w http.ResponseWriter, r *http.Request, hub *Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// clientID := c.Request.URL.Query().Get("id")
	clientID := "thisIsClientID"
	// clientID := r.URL.Query().Get("id")
	client := &Client{ID: clientID, Conn: conn, Send: make(chan []byte, 256), Hub: hub}
	// hub.Register <- client
	client.Hub.Register <- client

	go client.writePump()
	go client.readPump()
}
