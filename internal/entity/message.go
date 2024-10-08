package entity

import "time"

type Message struct {
	SenderUUID  string
	MessageUUID string
	Content     string
	CreatedAt   time.Time
}
type GetMessageDTO struct {
	MessageUUID string           `json:"message_uuid"`
	Content     string           `json:"content"`
	CreatedAt   time.Time        `json:"created_at"`
	User        UserProfileDTO   `json:"user"`
	Reaction    []GetReactionDTO `json:"reaction"`
}

type SeenStatus struct {
	UserUUID         string
	ConversationUUID string
}
type SeenStatusDTO struct {
	UserUUID         string
	ConversationUUID string
}
type GetSeenStatusDTO struct {
	UserProfileDTO
	SeenTimestamp string `json:"seen_timestamp"`
}

type MessageDTO struct {
	UserUUID    string
	MessageUUID string
}

type SearchMessageDTO struct {
	MessageUUID string         `json:"message_uuid"`
	Content     string         `json:"content"`
	CreatedAt   time.Time      `json:"created_at"`
	User        UserProfileDTO `json:"user"`
	Cursor      string         `json:"cursor"`
}
