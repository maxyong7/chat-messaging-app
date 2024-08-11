package entity

import "time"

type Conversation struct {
	SenderUUID       string
	ConversationUUID string
	MessageUUID      string
	Content          string
	CreatedAt        time.Time
}

type ConversationDTO struct {
	SenderUUID       string
	ConversationUUID string
	MessageUUID      string
	Content          string
	CreatedAt        time.Time
}

type ConversationList struct {
	ConversationUUID     *string     `json:"conversation_uuid"`
	Title                *string     `json:"title"`
	LastMessage          *string     `json:"last_message"`
	LastSentUser         UserProfile `json:"last_sent_user"`
	LastMessageCreatedAt *time.Time  `json:"last_message_created_at"`
	Type                 *string     `json:"type"`
}
