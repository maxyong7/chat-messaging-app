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
	User        UserInfoDTO      `json:"user"`
	Reaction    []GetReactionDTO `json:"reaction"`
}

type SeenStatus struct {
	UserInfo
	SeenTimestamp string `json:"seen_timestamp"`
}

type MessageDTO struct {
	UserUUID    string
	MessageUUID string
}
