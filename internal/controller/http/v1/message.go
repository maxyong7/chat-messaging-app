package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/boundary"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type messageRoute struct {
	t usecase.Message
	l logger.Interface
}

func newMessageRoute(handler *gin.RouterGroup, t usecase.Message, l logger.Interface) {
	route := &messageRoute{t, l}

	h := handler.Group("/message")
	{
		h.GET("/:conversation_uuid", route.getMessagesFromConversation)
	}
}

func (r *messageRoute) getMessagesFromConversation(c *gin.Context) {
	convUUID := c.Param("conversation_uuid")
	if convUUID == "" {
		errorResponse(c, http.StatusUnprocessableEntity, "missing conversation_uuid")
		return
	}

	cursor, err := queryParamCursor(c)
	if err != nil {
		r.l.Error(err, "http - v1 - getInbox - cursor validation error")
		handleCustomErrors(c, err)
	}
	limit := queryParamInt(c, "limit", 20)

	userId, err := getUserUUIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	requestParams := entity.RequestParams{
		Cursor: cursor,
		Limit:  limit,
		UserID: userId,
	}

	messages, err := r.t.GetMessagesFromConversation(c.Request.Context(), requestParams, convUUID)
	if err != nil {
		r.l.Error(err, "http - v1 - getContacts - GetContacts")
		handleCustomErrors(c, err)
		return
	}

	var encodedCursor string
	if len(messages) == limit {
		encodedCursor = encodeCursor(&messages[len(messages)-1].CreatedAt)
	}

	msgResp := boundary.MessageResponseModel{
		Data: boundary.MessageData{
			Messages: messages,
		},
		Pagination: boundary.Pagination{
			Cursor: encodedCursor,
			Limit:  limit,
		},
	}

	c.JSON(http.StatusOK, msgResp)
}
