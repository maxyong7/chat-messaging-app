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

const logLevelDebug = "debug"

func TestGetContacts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockContact(ctrl)
	mockLogger := logger.New(logLevelDebug)

	// Create a test instance of the contactRoute
	r := &contactRoute{t: mockUsecase, l: mockLogger}

	// Set up the request
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Mock getting user UUID from context
		userId := "some-uuid"
		c.Set("user_uuid", userId)
		c.Next()
	})

	router.GET("/contact", r.getContacts)
	t.Run("Success", func(t *testing.T) {

		// Mock the expected behavior
		userId := "some-uuid"
		mockUsecase.EXPECT().GetContacts(gomock.Any(), userId).Return([]entity.Contacts{}, nil)

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/contact", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the results
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Set up the request
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/contact", r.getContacts)

		// Perform the request without setting the user_uuid
		req, _ := http.NewRequest(http.MethodGet, "/contact", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the results
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Entity object failure - error calling GetContacts", func(t *testing.T) {

		// Mock the expected behavior
		userId := "some-uuid"
		mockUsecase.EXPECT().GetContacts(gomock.Any(), userId).Return([]entity.Contacts{}, errors.New("test_error"))

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/contact", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the results
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAddContact(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockContact(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &contactRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userId := "some-uuid"
		c.Set("user_uuid", userId)
		c.Next()
	})

	router.POST("/contact/:username/add", r.addContact)

	t.Run("Success", func(t *testing.T) {
		userId := "some-uuid"
		username := "testuser"
		mockUsecase.EXPECT().AddContact(gomock.Any(), username, userId).Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/contact/testuser/add", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("MissingUsername", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/contact//add", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/contact/:username/add", r.addContact)

		req, _ := http.NewRequest(http.MethodPost, "/contact/testuser/add", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Entity object failure - error calling AddContact", func(t *testing.T) {
		userId := "some-uuid"
		username := "testuser"
		mockUsecase.EXPECT().AddContact(gomock.Any(), username, userId).Return(errors.New("test_error"))

		req, _ := http.NewRequest(http.MethodPost, "/contact/testuser/add", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRemoveContact(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockContact(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &contactRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userId := "some-uuid"
		c.Set("user_uuid", userId)
		c.Next()
	})

	router.POST("/contact/:username/remove", r.removeContact)

	t.Run("Success", func(t *testing.T) {
		userId := "some-uuid"
		username := "testuser"
		mockUsecase.EXPECT().RemoveContact(gomock.Any(), username, userId).Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/contact/testuser/remove", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("MissingUsername", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/contact//remove", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/contact/:username/remove", r.removeContact)

		req, _ := http.NewRequest(http.MethodPost, "/contact/testuser/remove", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Entity object failure - error calling RemoveContact", func(t *testing.T) {
		userId := "some-uuid"
		username := "testuser"
		mockUsecase.EXPECT().RemoveContact(gomock.Any(), username, userId).Return(errors.New("test_error"))

		req, _ := http.NewRequest(http.MethodPost, "/contact/testuser/remove", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
func TestUpdateBlockContact(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockContact(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &contactRoute{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		userId := "some-uuid"
		c.Set("user_uuid", userId)
		c.Next()
	})

	router.PATCH("/contact/:username", r.updateBlockContact)

	t.Run("Success", func(t *testing.T) {
		userId := "some-uuid"
		username := "testuser"
		mockUsecase.EXPECT().UpdateBlockContact(gomock.Any(), username, userId, true).Return(nil)

		req, _ := http.NewRequest(http.MethodPatch, "/contact/testuser?block=true", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("InvalidBlockQuery", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/contact/testuser?block=invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PATCH("/contact/:username", r.updateBlockContact)

		req, _ := http.NewRequest(http.MethodPatch, "/contact/testuser?block=true", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Entity object failure - error calling UpdateBlockContact", func(t *testing.T) {
		userId := "some-uuid"
		username := "testuser"
		mockUsecase.EXPECT().UpdateBlockContact(gomock.Any(), username, userId, true).Return(errors.New("test_error"))

		req, _ := http.NewRequest(http.MethodPatch, "/contact/testuser?block=true", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
