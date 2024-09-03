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
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	mocks "github.com/maxyong7/chat-messaging-app/internal/usecase/mocks"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

func TestGetUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserProfile(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &userProfileRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userUUID := "some-uuid"
		c.Set("user_uuid", userUUID)
		c.Next()
	})

	router.GET("/profile", r.getUserProfile)

	t.Run("Success", func(t *testing.T) {
		userProfile := entity.UserProfile{
			UserUUID:  "some-uuid",
			FirstName: "John",
			LastName:  "Doe",
			Avatar:    "test_avatar",
		}
		mockUsecase.EXPECT().GetUserProfile(gomock.Any(), "some-uuid").Return(userProfile, nil)

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/profile", r.getUserProfile)

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("EntityObjectFailure", func(t *testing.T) {
		mockUsecase.EXPECT().GetUserProfile(gomock.Any(), "some-uuid").Return(entity.UserProfile{}, errors.New("test_error"))

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserProfile(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &userProfileRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userUUID := "some-uuid"
		c.Set("user_uuid", userUUID)
		c.Next()
	})

	router.PATCH("/profile", r.updateUserProfile)

	t.Run("Success", func(t *testing.T) {
		mockUsecase.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(nil)

		request := boundary.ProfileSettingScreen{
			FirstName: "Mary",
			LastName:  "Jane",
			Avatar:    "changed_avatar",
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPatch, "/profile", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PATCH("/profile", r.updateUserProfile)

		req, _ := http.NewRequest(http.MethodPatch, "/profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		requestBody := `{"name": "John Doe", "email": }` // Invalid JSON
		req, _ := http.NewRequest(http.MethodPatch, "/profile", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("EntityObjectFailure", func(t *testing.T) {
		mockUsecase.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(errors.New("test_error"))

		request := boundary.ProfileSettingScreen{
			FirstName: "Mary",
			LastName:  "Jane",
			Avatar:    "changed_avatar",
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPatch, "/profile", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
