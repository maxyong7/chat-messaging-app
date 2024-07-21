package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
)

type conversationRoutes struct {
	t usecase.Conversation
	l logger.Interface
}

func newConversationRoute(handler *gin.RouterGroup, t usecase.Conversation, l logger.Interface) {
	route := &conversationRoutes{t, l}
	hub := usecase.NewHub()
	go hub.Run()

	h := handler.Group("/conversation")
	{
		h.GET("/ws/:conversationId", route.ServeWsController(hub))
		// h.GET("/ws/:clientId", func(c *gin.Context) {

		// 	route.t.ServeWs(c, hub)
		// })
		// http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		// 	route.t.ServeWsWithRW(w, r, hub)
		// })
		// http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		// 	route.t.ServeWsWithRW(w, r, hub)
		// })
		// h.GET("/ws", route.ServeWsController)
	}
}

func (r *conversationRoutes) ServeWsController(hub *usecase.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := getUserUUIDFromContext(c)
		if err != nil {
			errorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		// hub := usecase.NewHub()
		// go hub.Run()

		r.t.ServeWs(c, hub, userId)
	}
}

// func main() {
// 	server := NewServer()
// 	http.Handle("/ws".websocket.Handler(server.handleWS))
// 	http.ListenAndServe(":3000,nil")
// }
