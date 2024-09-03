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

func TestGetMessagesFromConversation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockMessage(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &messageRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userUUID := "some-uuid"
		c.Set("user_uuid", userUUID)
		c.Next()
	})

	router.GET("/message/:conversation_uuid", r.getMessagesFromConversation)

	t.Run("Success", func(t *testing.T) {
		convUUID := "conv-uuid"
		messages := []entity.GetMessageDTO{
			{Content: "Message 1"},
			{Content: "Message 2"},
		}
		mockUsecase.EXPECT().GetMessagesFromConversation(gomock.Any(), gomock.Any(), convUUID).Return(messages, nil)
		mockUsecase.EXPECT().UpdateSeenStatus(gomock.Any(), gomock.Any()).Return(nil)

		req, _ := http.NewRequest(http.MethodGet, "/message/conv-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/message/:conversation_uuid", r.getMessagesFromConversation)

		req, _ := http.NewRequest(http.MethodGet, "/message/conv-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("EntityObjectFailure", func(t *testing.T) {
		convUUID := "conv-uuid"
		mockUsecase.EXPECT().GetMessagesFromConversation(gomock.Any(), gomock.Any(), convUUID).Return(nil, errors.New("test_error"))

		req, _ := http.NewRequest(http.MethodGet, "/message/conv-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestSearchMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockMessage(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &messageRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userUUID := "some-uuid"
		c.Set("user_uuid", userUUID)
		c.Next()
	})

	router.GET("/message/:conversation_uuid/search", r.searchMessage)

	t.Run("Success", func(t *testing.T) {
		convUUID := "conv-uuid"
		messages := []entity.SearchMessageDTO{
			{Content: "Message containing keyword"},
		}
		mockUsecase.EXPECT().SearchMessage(gomock.Any(), "keyword", convUUID).Return(messages, nil)

		req, _ := http.NewRequest(http.MethodGet, "/message/conv-uuid/search?keyword=keyword", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("MissingConversationUUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/message//search?keyword=keyword", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("MissingKeyword", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/message/conv-uuid/search", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("EntityObjectFailure", func(t *testing.T) {
		convUUID := "conv-uuid"
		mockUsecase.EXPECT().SearchMessage(gomock.Any(), "keyword", convUUID).Return(nil, errors.New("test_error"))

		req, _ := http.NewRequest(http.MethodGet, "/message/conv-uuid/search?keyword=keyword", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetMessageStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockMessage(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &messageRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userUUID := "some-uuid"
		c.Set("user_uuid", userUUID)
		c.Next()
	})

	router.GET("/message/status/:message_uuid", r.getMessageStatus)

	t.Run("Success", func(t *testing.T) {
		msgUUID := "msg-uuid"
		seenStatus := []entity.GetSeenStatusDTO{
			{
				UserProfileDTO: entity.UserProfileDTO{
					UserUUID:  "test_useruuid",
					FirstName: "test_firstname",
					LastName:  "test_lastname",
					Avatar:    "test_avatar",
				},
				SeenTimestamp: "2024-01-01T00:00:00Z",
			},
		}
		mockUsecase.EXPECT().GetSeenStatus(gomock.Any(), msgUUID).Return(seenStatus, nil)

		req, _ := http.NewRequest(http.MethodGet, "/message/status/msg-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("EntityObjectFailure", func(t *testing.T) {
		msgUUID := "msg-uuid"
		mockUsecase.EXPECT().GetSeenStatus(gomock.Any(), msgUUID).Return([]entity.GetSeenStatusDTO{}, errors.New("test_error"))

		req, _ := http.NewRequest(http.MethodGet, "/message/status/msg-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
