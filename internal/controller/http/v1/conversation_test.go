package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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

func TestServeWsController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConversationUsecase := mocks.NewMockConversation(ctrl)
	mockUserProfileUsecase := mocks.NewMockUserProfile(ctrl)
	mockLogger := logger.New(logLevelDebug)

	// Create a test instance of the conversationRoutes
	r := &conversationRoutes{conv: mockConversationUsecase, up: mockUserProfileUsecase, l: mockLogger}

	// Set up the request
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userId := "some-uuid"
		c.Set("user_uuid", userId)
		c.Next()
	})
	hub := NewHub()
	go hub.Run()

	router.GET("/conversation/ws/:conversationId", r.serveWsController(hub))

	t.Run("Unauthorized", func(t *testing.T) {
		// Set up the request without user_uuid context
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/conversation/ws/:conversationId", r.serveWsController(hub))

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/conversation/ws/some-conversation-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the results
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Success - WebSocket Connection", func(t *testing.T) {
		// Mock the expected behavior
		userId := "some-uuid"
		mockUserInfo := entity.UserProfile{UserUUID: userId}

		// Set up the context to include user_uuid
		router.Use(func(c *gin.Context) {
			c.Set("user_uuid", userId)
			c.Next()
		})

		// Mock the GetUserProfile method
		mockUserProfileUsecase.EXPECT().GetUserProfile(gomock.Any(), userId).Return(mockUserInfo, nil)

		// Create a request for the WebSocket connection
		req, _ := http.NewRequest(http.MethodGet, "/conversation/ws/some-conversation-id", nil)
		w := httptest.NewRecorder()

		// Set required WebSocket headers
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Sec-WebSocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "x3JJHMbDL1EzLkh9GBhXDw==")

		// Serve the request
		router.ServeHTTP(w, req)

		// Assert the response for a successful WebSocket upgrade
		assert.Equal(t, http.StatusSwitchingProtocols, w.Code)
	})
}
