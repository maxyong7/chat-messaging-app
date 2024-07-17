package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type contactRoute struct {
	t usecase.Contact
	l logger.Interface
}

func newContactRoute(handler *gin.RouterGroup, t usecase.Contact, l logger.Interface) {
	route := &contactRoute{t, l}

	h := handler.Group("/contact")
	{
		h.GET("", route.getContacts)
		h.POST("/:username/add", route.addContact)
		h.POST("/:username/remove", route.removeContact)
		h.PATCH("/:username", route.updateBlockContact)
	}
}

func (r *contactRoute) getContacts(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	contacts, err := r.t.GetContacts(c.Request.Context(), userId)
	if err != nil {
		r.l.Error(err, "http - v1 - getContacts - GetContacts")
		handleCustomErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, contacts)
}

func (r *contactRoute) addContact(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	contactUserName := c.Param("username")
	if contactUserName == "" {
		errorResponse(c, http.StatusUnprocessableEntity, "missing username in parameter")
		return
	}

	err = r.t.AddContact(c.Request.Context(), contactUserName, userId)
	if err != nil {
		r.l.Error(err, "http - v1 - addContact - AddContacts")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}

func (r *contactRoute) removeContact(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	contactUserName := c.Param("username")
	if contactUserName == "" {
		errorResponse(c, http.StatusUnprocessableEntity, "missing username in parameter")
		return
	}
	err = r.t.RemoveContact(c.Request.Context(), contactUserName, userId)
	if err != nil {
		r.l.Error(err, "http - v1 - addContact - AddContacts")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}

func (r *contactRoute) updateBlockContact(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	contactUserName := c.Param("username")
	if contactUserName == "" {
		errorResponse(c, http.StatusUnprocessableEntity, "missing username in parameter")
		return
	}
	blockString, ok := c.GetQuery("block")
	if !ok {
		errorResponse(c, http.StatusUnprocessableEntity, "block query missing")
		return
	}
	var block bool
	switch strings.ToLower(blockString) {
	case "true":
		block = true
	case "false":
		block = false
	default:
		errorResponse(c, http.StatusUnprocessableEntity, "block query missing")
		return
	}

	err = r.t.UpdateBlockContact(c.Request.Context(), contactUserName, userId, block)
	if err != nil {
		r.l.Error(err, "http - v1 - updateBlockContact - UpdateBlockContact")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
