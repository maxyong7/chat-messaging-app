package entity

import "time"

type InboxResponse struct {
	Data       Data       `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Data struct {
	Conversations []Conversations `json:"conversations"`
}

type Conversations struct {
	Title                string               `json:"title"`
	LastMessage          string               `json:"last_message"`
	LastSentUser         ConversationUserInfo `json:"last_sent_user"`
	LastMessageCreatedAt time.Time            `json:"last_message_created_at"`
	Type                 string               `json:"type"`
}

type ConversationUserInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

type ConversationDTO struct {
	ConversationUUID     string    `db:"conversation_uuid"`
	LastMessage          string    `db:"last_message"`
	LastSentUserUUID     string    `db:"last_sent_user_uuid"`
	Title                string    `db:"title"`
	LastMessageCreatedAt time.Time `db:"last_message_created_at"`
	ConversationType     string    `db:"conversation_type"`
}