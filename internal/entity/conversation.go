package entity

import (
	"time"
)

type Conversation struct {
	SenderUUID       string
	ConversationUUID string
	ConversationMessage
}

type ConversationMessage struct {
	MessageUUID string    `json:"message_uuid,omitempty"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}
