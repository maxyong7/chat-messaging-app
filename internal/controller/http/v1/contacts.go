package v1

import (
	"net/http"

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
	err = r.t.RemoveContact(c.Request.Context(), contactUserName, userId)
	if err != nil {
		r.l.Error(err, "http - v1 - addContact - AddContacts")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
