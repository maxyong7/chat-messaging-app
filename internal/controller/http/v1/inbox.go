package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type inboxRoute struct {
	t usecase.Inbox
	l logger.Interface
}

func newInboxRoute(handler *gin.RouterGroup, t usecase.Inbox, l logger.Interface) {
	route := &inboxRoute{t, l}

	h := handler.Group("/inbox")
	{
		h.GET("", route.getInbox)
		// http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		// 	route.t.ServeWsWithRW(w, r, hub)
		// })
		// http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		// 	route.t.ServeWsWithRW(w, r, hub)
		// })
		// h.GET("/ws", route.ServeWsController)
	}
}

func (r *inboxRoute) getInbox(c *gin.Context) {
	cursor, err := queryParamCursor(c)
	if err != nil {
		r.l.Error(err, "http - v1 - getInbox - cursor validation error")
		handleCustomErrors(c, err)
	}
	limit := queryParamInt(c, "limit", 20)

	userId, err := getUserIDFromContext(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	requestParams := entity.RequestParams{
		Cursor: cursor,
		Limit:  limit,
		UserID: userId,
	}

	inboxResponse, err := r.t.GetInbox(c.Request.Context(), requestParams)
	if err != nil {
		r.l.Error(err, "http - v1 - getInbox - GetInbox")
		handleCustomErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, inboxResponse)
}
