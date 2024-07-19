package entity

import "time"

type MessageResponse struct {
	Data       MessageData `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type MessageData struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	MessageUUID string    `json:"message_uuid"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	User        UserInfo
	Reaction    []Reaction
}

type SeenStatus struct {
	UserInfo
	SeenTimestamp string `json:"seen_timestamp"`
}
