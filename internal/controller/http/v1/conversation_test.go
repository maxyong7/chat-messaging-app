package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"github.com/maxyong7/chat-messaging-app/internal/boundary"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	mocks "github.com/maxyong7/chat-messaging-app/internal/usecase/mocks"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

func TestGetConversations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockConversation(ctrl)
	mockLogger := logger.New(logLevelDebug)

	// Create a test instance of the conversationRoutes
	r := &conversationRoutes{conv: mockUsecase, l: mockLogger}

	// Set up the request
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Mock getting user UUID from context
		userId := "some-uuid"
		c.Set("user_uuid", userId)
		c.Next()
	})

	router.GET("/conversation", r.getConversations)
	t.Run("Success", func(t *testing.T) {
		// Mock the expected behavior
		mockConversations := []entity.ConversationList{}
		mockUsecase.EXPECT().GetConversationList(gomock.Any(), gomock.Any()).Return(mockConversations, nil)

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/conversation", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the results
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Set up the request without user_uuid context
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/conversation", r.getConversations)

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/conversation", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the results
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Entity object failure - error calling GetConversationList", func(t *testing.T) {
		// Mock the expected behavior
		mockUsecase.EXPECT().GetConversationList(gomock.Any(), gomock.Any()).Return(nil, errors.New("test_error"))

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/conversation", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the results
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// func TestServeWsController(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockUserProfileUsecase := mocks.NewMockUserProfile(ctrl)
// 	mockLogger := logger.New(logLevelDebug)

// 	r := &conversationRoutes{up: mockUserProfileUsecase, l: mockLogger}

// 	gin.SetMode(gin.TestMode)
// 	router := gin.New()
// 	hub := NewHub()
// 	go hub.Run()

// 	router.Use(func(c *gin.Context) {
// 		userUUID := "some-uuid"
// 		c.Set("user_uuid", userUUID)
// 		c.Next()
// 	})

// 	router.GET("/conversation/ws/:conversationId", r.serveWsController(hub))

// 	t.Run("Success", func(t *testing.T) {
// 		mockUserProfileUsecase.EXPECT().GetUserProfile(gomock.Any(), "some-uuid").Return(entity.UserProfile{
// 			UserUUID: "some-uuid",
// 		}, nil)

// 		server := httptest.NewServer(router)
// 		defer server.Close()

// 		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/conversation/ws/conv-uuid"

// 		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 		assert.NoError(t, err)
// 		defer ws.Close()

// 		assert.Equal(t, websocket.StateConnected, ws.UnderlyingConn().State())
// 	})

// 	t.Run("Unauthorized", func(t *testing.T) {
// 		router := gin.New()
// 		router.GET("/conversation/ws/:conversationId", r.serveWsController(hub))

// 		server := httptest.NewServer(router)
// 		defer server.Close()

// 		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/conversation/ws/conv-uuid"

// 		_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 		assert.Error(t, err)
// 	})
// }

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func TestHandleConversation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConvUsecase := mocks.NewMockConversation(ctrl)
	mockMsgUsecase := mocks.NewMockMessage(ctrl)
	mockReactionUsecase := mocks.NewMockReaction(ctrl)
	mockLogger := logger.New(logLevelDebug)

	hub := NewHub()
	go hub.Run()

	r := &conversationRoutes{
		conv:     mockConvUsecase,
		msg:      mockMsgUsecase,
		reaction: mockReactionUsecase,
		l:        mockLogger,
	}

	userInfo := entity.UserProfile{
		UserUUID: "some-uuid",
	}
	// userInfo2 := entity.UserProfile{
	// 	UserUUID: "some-uuid2",
	// }

	// router := gin.New()
	// router.Use(func(c *gin.Context) {
	// 	userUUID := "some-uuid"
	// 	c.Set("user_uuid", userUUID)
	// 	c.Next()
	// })

	// router.GET("/conversation/ws/:conversationId", r.serveWsController(hub))

	// server := httptest.NewServer(router)

	server := httptest.NewServer(http.HandlerFunc(echo))
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, _, _ := websocket.DefaultDialer.Dial(u, nil)
	defer conn.Close()
	conn2, _, _ := websocket.DefaultDialer.Dial(u, nil)
	defer conn.Close()

	client := NewClient("conv-uuid", userInfo, conn, hub, r)
	client2 := NewClient("conv-uuid", userInfo, conn2, hub, r)
	// client2 := NewClient("conv-uuid", userInfo2, conn, hub, r)
	t.Run("SendMessageSuccess", func(t *testing.T) {
		// conn, _, _ := websocket.DefaultDialer.Dial("ws://localhost:8080", nil)
		// defer conn.Close()
		// conn2, _, _ := websocket.DefaultDialer.Dial("ws://localhost:8080", nil)
		// defer conn2.Close()

		client := NewClient("conv-uuid", userInfo, conn, hub, r)
		hub.Register <- client
		hub.Register <- client2
		msgContent := "Hello, World!"
		msgData, _ := json.Marshal(boundary.ChatInterface{Content: msgContent})
		convReq := boundary.ConversationRequestModel{
			MessageType: sendMessageType,
			Data:        msgData,
		}

		// var wg sync.WaitGroup
		// ready := make(chan struct{}) // Channel to signal that the broadcast goroutine is ready
		// wg.Add(1)

		mockConvUsecase.EXPECT().StoreConversationAndMessage(gomock.Any(), gomock.Any()).Return(nil)

		// Start a goroutine that listens for the broadcast message before handling the conversation
		// go func() {
		// 	defer wg.Done()
		// 	// Signal that the goroutine is ready to listen
		// 	close(ready)

		// 	_, p, err := client.Conn.ReadMessage()
		// 	if err != nil {
		// 		t.Fatalf("%v", err)
		// 	}
		// 	msg := string(p)

		// 	assert.Equal(t, msgContent, msg)
		// }()

		// // Ensure the listener goroutine is ready
		// <-ready
		// time.Sleep(100 * time.Millisecond) // Optional sleep to ensure the goroutine is fully ready

		// Trigger the handleConversation method
		client.handleConversation(convReq, userInfo)

		// time.Sleep(1000 * time.Millisecond) // Optional sleep to ensure the goroutine is fully ready

		// for i := 0; i < 10; i++ {
		// 	if err := conn.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
		// 		t.Fatalf("%v", err)
		// 	}
		// 	_, p, err := conn.ReadMessage()
		// 	if err != nil {
		// 		t.Fatalf("%v", err)
		// 	}
		// 	if string(p) != "hello" {
		// 		t.Fatalf("bad message")
		// 	}
		// }
		// _, p, err := conn2.ReadMessage()
		// _, p, err := conn.ReadMessage()
		_, p, err := client2.Conn.ReadMessage()
		// _, p, err := client.Conn.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}
		msg := string(p)

		assert.Equal(t, msgContent, msg)
		// assert.Equal(t, sendMessageType, msg.MessageType)
		// assert.Equal(t, msgContent, msg.Data.SendMessageResponseData.Content)

		// Wait for the message to be processed
		// wg.Wait()

		// if client.Conn != nil {
		// 	err := client.Conn.Close()
		// 	if err != nil {
		// 		t.Errorf("Error closing WebSocket connection: %v", err)
		// 	}
		// }
	})

	t.Run("DeleteMessageSuccess", func(t *testing.T) {
		deleteReq := boundary.DeleteMessageRequest{
			MessageUUID: "msg-uuid",
		}
		msgData, _ := json.Marshal(deleteReq)
		convReq := boundary.ConversationRequestModel{
			MessageType: deleteMessageType,
			Data:        msgData,
		}

		mockMsgUsecase.EXPECT().DeleteMessage(gomock.Any(), gomock.Any()).Return(true, nil)

		client.handleConversation(convReq, userInfo)

		msg := <-client.hub.Broadcast
		assert.Equal(t, deleteMessageType, msg.MessageType)
		assert.Equal(t, "msg-uuid", msg.Data.MessageDeletionConfirmation.MessageUUID)
	})

	t.Run("AddReactionSuccess", func(t *testing.T) {
		reactionReq := boundary.MessageReactionMenu{
			MessageUUID:  "msg-uuid",
			ReactionType: "like",
		}
		msgData, _ := json.Marshal(reactionReq)
		convReq := boundary.ConversationRequestModel{
			MessageType: addReactionMessageType,
			Data:        msgData,
		}

		mockReactionUsecase.EXPECT().StoreReaction(gomock.Any(), gomock.Any()).Return(nil)

		client.handleConversation(convReq, userInfo)

		msg := <-client.hub.Broadcast
		assert.Equal(t, addReactionMessageType, msg.MessageType)
		assert.Equal(t, "like", msg.Data.ReactionResponseData.Reaction)
	})

	t.Run("RemoveReactionSuccess", func(t *testing.T) {
		removeReactionReq := boundary.RemoveReactionRequest{
			MessageUUID: "msg-uuid",
		}
		msgData, _ := json.Marshal(removeReactionReq)
		convReq := boundary.ConversationRequestModel{
			MessageType: removeReactionMessageType,
			Data:        msgData,
		}

		mockReactionUsecase.EXPECT().RemoveReaction(gomock.Any(), gomock.Any()).Return(nil)

		client.handleConversation(convReq, userInfo)

		msg := <-client.hub.Broadcast
		assert.Equal(t, removeReactionMessageType, msg.MessageType)
		assert.Equal(t, "msg-uuid", msg.Data.ReactionResponseData.MessageUUID)
	})

	t.Run("HandleErrorForInvalidMessage", func(t *testing.T) {
		convReq := boundary.ConversationRequestModel{
			MessageType: "invalid_type",
			Data:        nil,
		}

		client.handleConversation(convReq, userInfo)

		msg := <-client.hub.Broadcast
		assert.Equal(t, errorMessageType, msg.MessageType)
	})
}
