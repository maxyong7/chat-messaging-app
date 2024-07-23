package entity

import "time"

type GetMessageDTO struct {
	MessageUUID string    `json:"message_uuid"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	User        UserInfo
	Reaction    []Reaction
}

type Message struct {
	MessageUUID string    `json:"message_uuid"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

type SeenStatus struct {
	UserInfo
	SeenTimestamp string `json:"seen_timestamp"`
}
