package entity

import "time"

type GetMessageDTO struct {
	MessageUUID string    `json:"message_uuid"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	User        UserInfo
	Reaction    []GetReaction
}

type Message struct {
	SenderUUID  string
	MessageUUID string
	Content     string
	CreatedAt   time.Time
}

type SeenStatus struct {
	UserInfo
	SeenTimestamp string `json:"seen_timestamp"`
}

type MessageDTO struct {
	UserUUID    string
	MessageUUID string
}
