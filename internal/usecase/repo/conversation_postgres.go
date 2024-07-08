package repo

import (
	"github.com/maxyong7/chat-messaging-app/pkg/postgres"
)

// ConversationRepo -.
type ConversationRepo struct {
	*postgres.Postgres
}

// New -.
func NewConversation(pg *postgres.Postgres) *ConversationRepo {
	return &ConversationRepo{pg}
}
