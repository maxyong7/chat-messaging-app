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

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUser(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &userRoutes{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/user/login", r.loginUser)

	t.Run("Success", func(t *testing.T) {
		mockUsecase.EXPECT().VerifyCredentials(gomock.Any(), gomock.Any()).Return("some-uuid", true, nil)

		request := boundary.LoginForm{
			Username: "testjohndoe",
			Password: "password123",
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodPost, "/user/login", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		requestBody := `{"username": "testjohndoe", "password": }` // Invalid JSON
		req, _ := http.NewRequest(http.MethodPost, "/user/login", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		request := boundary.LoginForm{
			Username: "testjohndoe",
			Password: "wrongpassword",
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		mockUsecase.EXPECT().VerifyCredentials(gomock.Any(), gomock.Any()).Return("", false, entity.ErrIncorrectPassword)
		req, _ := http.NewRequest(http.MethodPost, "/user/login", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("EntityObjectFailure", func(t *testing.T) {
		request := boundary.LoginForm{
			Username: "testjohndoe",
			Password: "wrongpassword",
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		mockUsecase.EXPECT().VerifyCredentials(gomock.Any(), gomock.Any()).Return("", false, errors.New("test_error"))
		req, _ := http.NewRequest(http.MethodPost, "/user/login", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUser(ctrl)
	mockLogger := logger.New(logLevelDebug)

	r := &userRoutes{t: mockUsecase, l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/user/register", r.registerUser)

	t.Run("Success", func(t *testing.T) {
		request := boundary.RegistrationForm{
			Username:  "testjohndoe",
			Password:  "password123",
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Avatar:    "testavatar",
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		mockUsecase.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(nil)
		req, _ := http.NewRequest(http.MethodPost, "/user/register", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		requestBody := `{"email": "test@example.com", "password": }` // Invalid JSON
		req, _ := http.NewRequest(http.MethodPost, "/user/register", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("EntityObjectFailure", func(t *testing.T) {
		request := boundary.RegistrationForm{
			Username:  "testjohndoe",
			Password:  "password123",
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Avatar:    "testavatar",
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(request)
		if err != nil {
			log.Fatal(err)
		}

		mockUsecase.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(errors.New("test_error"))
		req, _ := http.NewRequest(http.MethodPost, "/user/register", &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLogoutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.New(logLevelDebug)

	r := &userRoutes{l: mockLogger}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/user/logout", r.logoutUser)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/user/logout", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	// t.Run("TokenCreationFailure", func(t *testing.T) {
	// 	// Mock the failure in token creation
	// 	createToken = func(uuid string, hours int) (string, error) {
	// 		return "", errors.New("token creation error")
	// 	}

	// 	req, _ := http.NewRequest(http.MethodPost, "/user/logout", nil)
	// 	w := httptest.NewRecorder()

	// 	router.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// })
}
