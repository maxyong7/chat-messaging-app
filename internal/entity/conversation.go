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
