package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func handleCustomErrors(c *gin.Context, err error) {
	switch err {
	case entity.ErrUserAlreadyExists, entity.ErrContactAlreadyExists:
		errorResponse(c, http.StatusConflict, err.Error())
	case entity.ErrUserNameNotFound, entity.ErrContactDoesNotExists, entity.ErrUserNotFound:
		errorResponse(c, http.StatusNotFound, err.Error())
	case entity.ErrIncorrectPassword:
		errorResponse(c, http.StatusUnauthorized, err.Error())
	default:
		errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}
