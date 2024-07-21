package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/boundary"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type webserverRoutes struct {
	t usecase.Verification
	l logger.Interface
}

func newUserVerificationRoute(handler *gin.RouterGroup, t usecase.Verification, l logger.Interface) {
	r := &webserverRoutes{t, l}

	h := handler.Group("/user")
	{
		h.POST("/verify", r.verifyCredentials)
		h.POST("/register", r.registerUser)
	}
}

func (r *webserverRoutes) verifyCredentials(c *gin.Context) {
	var request boundary.VerifyUserRequestModel
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - verifyCredentials")
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}

	userUuid, isValid, err := r.t.VerifyCredentials(c.Request.Context(), request.ToVerifyCredentials())
	if err != nil {
		r.l.Error(err, "http - v1 - verifyCredentials")
		handleCustomErrors(c, err)
		return
	}

	if isValid && userUuid != "" {
		token, err := createToken(userUuid)
		if err != nil {
			r.l.Error(err, "http - v1 - createToken")
			handleCustomErrors(c, err)
			return
		}
		c.JSON(http.StatusOK, boundary.VerifyUserResponseModel{
			Token: token,
		})
		return
	}

	handleCustomErrors(c, err)
}

func (r *webserverRoutes) registerUser(c *gin.Context) {
	var request boundary.RegisterUserRequestModel
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - registerUser")
		errorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}

	err := r.t.RegisterUser(c.Request.Context(), request.ToUserRegistration())
	if err != nil {
		r.l.Error(err, "http - v1 - registerUser")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}
