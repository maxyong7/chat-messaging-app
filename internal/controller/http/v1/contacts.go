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
		h.POST("/:username/add", route.addContact)
		// http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		// 	route.t.ServeWsWithRW(w, r, hub)
		// })
		// http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		// 	route.t.ServeWsWithRW(w, r, hub)
		// })
		// h.GET("/ws", route.ServeWsController)
	}
}

func (r *contactRoute) addContact(c *gin.Context) {
	userId, err := getUserIDFromContext(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	contactUserName := c.Param("username")
	r.t.AddContacts(c.Request.Context(), contactUserName, userId)
	if err != nil {
		r.l.Error(err, "http - v1 - getInbox - GetInbox")
		handleCustomErrors(c, err)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}
