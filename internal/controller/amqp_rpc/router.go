package amqprpc

import (
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter -.
func NewRouter(t usecase.Translation) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newTranslationRoutes(routes, t)
	}

	return routes
}
