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
