package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/maxyong7/chat-messaging-app/internal/boundary"
	mocks "github.com/maxyong7/chat-messaging-app/internal/usecase/mocks"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

func TestCreateGroupChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockGroupChat(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &groupChatRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userUUID := "some-uuid"
		c.Set("user_uuid", userUUID)
		c.Next()
	})

	router.POST("/groupchat/create", r.createGroupChat)

	t.Run("Success", func(t *testing.T) {
		mockUsecase.EXPECT().CreateGroupChat(gomock.Any(), gomock.Any()).Return(nil)

		request := boundary.GroupChatCreationForm{
			Title:            "Test Group",
			ConversationUUID: "test_conversation_uuid",
			Participants: []boundary.ParticipantsForm{{
				Username:        "test_user",
				ParticipantUUID: "user_uuid1234",
			}},
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPost, "/groupchat/create", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		router := gin.New()
		router.POST("/groupchat/create", r.createGroupChat)

		request := boundary.GroupChatCreationForm{
			Title:            "Test Group",
			ConversationUUID: "test_conversation_uuid",
			Participants: []boundary.ParticipantsForm{{
				Username:        "test_user",
				ParticipantUUID: "user_uuid1234",
			}},
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPost, "/groupchat/create", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		reqBody := `{"invalid_json"}`
		req, _ := http.NewRequest(http.MethodPost, "/groupchat/create", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("CreateGroupChatFailure", func(t *testing.T) {
		mockUsecase.EXPECT().CreateGroupChat(gomock.Any(), gomock.Any()).Return(errors.New("test_error"))

		request := boundary.GroupChatCreationForm{
			Title:            "Test Group",
			ConversationUUID: "test_conversation_uuid",
			Participants: []boundary.ParticipantsForm{{
				Username:        "test_user",
				ParticipantUUID: "user_uuid1234",
			}},
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPost, "/groupchat/create", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAddParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockGroupChat(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &groupChatRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userUUID := "some-uuid"
		c.Set("user_uuid", userUUID)
		c.Next()
	})

	router.POST("/groupchat/add", r.addParticipant)

	t.Run("Success", func(t *testing.T) {
		mockUsecase.EXPECT().AddParticipant(gomock.Any(), gomock.Any()).Return(nil)
		request := boundary.GroupChatCreationForm{
			Title:            "Test Group",
			ConversationUUID: "test_conversation_uuid",
			Participants: []boundary.ParticipantsForm{{
				Username:        "test_user",
				ParticipantUUID: "user_uuid1234",
			}},
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPost, "/groupchat/add", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		router := gin.New()
		router.POST("/groupchat/add", r.addParticipant)

		request := boundary.GroupChatCreationForm{
			Title:            "Test Group",
			ConversationUUID: "test_conversation_uuid",
			Participants: []boundary.ParticipantsForm{{
				Username:        "test_user",
				ParticipantUUID: "user_uuid1234",
			}},
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPost, "/groupchat/add", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		reqBody := `{"invalid_json"}`
		req, _ := http.NewRequest(http.MethodPost, "/groupchat/add", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("AddParticipantFailure", func(t *testing.T) {
		mockUsecase.EXPECT().AddParticipant(gomock.Any(), gomock.Any()).Return(errors.New("test_error"))

		request := boundary.GroupChatCreationForm{
			Title:            "Test Group",
			ConversationUUID: "test_conversation_uuid",
			Participants: []boundary.ParticipantsForm{{
				Username:        "test_user",
				ParticipantUUID: "user_uuid1234",
			}},
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPost, "/groupchat/add", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRemoveParticipant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockGroupChat(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &groupChatRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userUUID := "some-uuid"
		c.Set("user_uuid", userUUID)
		c.Next()
	})

	router.POST("/groupchat/remove", r.removeParticipant)

	t.Run("Success", func(t *testing.T) {
		mockUsecase.EXPECT().RemoveParticipant(gomock.Any(), gomock.Any()).Return(nil)
		request := boundary.GroupChatCreationForm{
			Title:            "Test Group",
			ConversationUUID: "test_conversation_uuid",
			Participants: []boundary.ParticipantsForm{{
				Username:        "test_user",
				ParticipantUUID: "user_uuid1234",
			}},
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPost, "/groupchat/remove", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		router := gin.New()
		router.POST("/groupchat/remove", r.removeParticipant)

		request := boundary.GroupChatCreationForm{
			Title:            "Test Group",
			ConversationUUID: "test_conversation_uuid",
			Participants: []boundary.ParticipantsForm{{
				Username:        "test_user",
				ParticipantUUID: "user_uuid1234",
			}},
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPost, "/groupchat/remove", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		reqBody := `{"invalid_json"}`
		req, _ := http.NewRequest(http.MethodPost, "/groupchat/remove", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("RemoveParticipantFailure", func(t *testing.T) {
		mockUsecase.EXPECT().RemoveParticipant(gomock.Any(), gomock.Any()).Return(errors.New("test_error"))

		request := boundary.GroupChatCreationForm{
			Title:            "Test Group",
			ConversationUUID: "test_conversation_uuid",
			Participants: []boundary.ParticipantsForm{{
				Username:        "test_user",
				ParticipantUUID: "user_uuid1234",
			}},
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPost, "/groupchat/remove", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateGroupTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockGroupChat(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &groupChatRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userUUID := "some-uuid"
		c.Set("user_uuid", userUUID)
		c.Next()
	})

	router.PATCH("/groupchat/title", r.updateGroupTitle)

	t.Run("Success", func(t *testing.T) {
		request := boundary.GroupChatCreationForm{
			Title:            "Updated Title",
			ConversationUUID: "test_conversation_uuid",
		}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		mockUsecase.EXPECT().UpdateGroupTitle(gomock.Any(), gomock.Any()).Return(nil)

		req, _ := http.NewRequest(http.MethodPatch, "/groupchat/title", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		router := gin.New()
		router.PATCH("/groupchat/title", r.updateGroupTitle)

		request := boundary.GroupChatCreationForm{
			Title:            "Updated Title",
			ConversationUUID: "test_conversation_uuid",
		}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPatch, "/groupchat/title", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		reqBody := `{"invalid_json"}`
		req, _ := http.NewRequest(http.MethodPatch, "/groupchat/title", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UpdateGroupTitleFailure", func(t *testing.T) {
		mockUsecase.EXPECT().UpdateGroupTitle(gomock.Any(), gomock.Any()).Return(errors.New("test_error"))

		request := boundary.GroupChatCreationForm{
			Title:            "Updated Title",
			ConversationUUID: "test_conversation_uuid",
		}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodPatch, "/groupchat/title", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
