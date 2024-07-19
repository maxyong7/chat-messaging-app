package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type groupChatRoute struct {
	t usecase.GroupChat
	l logger.Interface
}

func newGroupChatRoute(handler *gin.RouterGroup, t usecase.GroupChat, l logger.Interface) {
	route := &groupChatRoute{t, l}

	h := handler.Group("/groupchat")
	{
		h.POST("/create", route.createGroupChat)
		h.POST("/add", route.addParticipant)
		h.POST("/remove", route.removeParticipant)
		h.PATCH("/title", route.updateGroupTitle)
	}
}

func (r *groupChatRoute) createGroupChat(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var request entity.GroupChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - createGroupChat")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	request.UserUUID = userId

	err = r.t.CreateGroupChat(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err, "http - v1 - createGroupChat - CreateGroupChat")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}

func (r *groupChatRoute) addParticipant(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var request entity.GroupChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - addParticipant")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	request.UserUUID = userId

	err = r.t.AddParticipant(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err, "http - v1 - addParticipant - AddParticipant")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func (r *groupChatRoute) removeParticipant(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var request entity.GroupChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - removeParticipant")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	request.UserUUID = userId

	err = r.t.RemoveParticipant(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err, "http - v1 - removeParticipant - RemoveParticipant")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func (r *groupChatRoute) updateGroupTitle(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var request entity.GroupChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - updateGroupTitle")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	request.UserUUID = userId

	err = r.t.UpdateGroupTitle(c.Request.Context(), request)
	if err != nil {
		r.l.Error(err, "http - v1 - updateGroupTitle - UpdateGroupTitle")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
