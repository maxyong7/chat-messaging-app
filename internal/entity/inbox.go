package entity

import (
	"time"
)

type Conversations struct {
	Title                *string    `json:"title"`
	LastMessage          *string    `json:"last_message"`
	LastSentUser         UserInfo   `json:"last_sent_user"`
	LastMessageCreatedAt *time.Time `json:"last_message_created_at"`
	Type                 *string    `json:"type"`
}

type ConversationDTO struct {
	ConversationUUID     string    `db:"conversation_uuid"`
	LastMessage          string    `db:"last_message"`
	LastSentUserUUID     string    `db:"last_sent_user_uuid"`
	Title                string    `db:"title"`
	LastMessageCreatedAt time.Time `db:"last_message_created_at"`
	ConversationType     string    `db:"conversation_type"`
}
