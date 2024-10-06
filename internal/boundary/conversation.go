package boundary

import (
	"encoding/json"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type ConversationRequestModel struct {
	MessageType string          `json:"message_type" binding:"required"`
	Data        json.RawMessage `json:"data" binding:"required"`
}

// type ConversationRequestData struct {
// 	SendMessageRequest
// 	DeleteMessageRequest
// 	AddReactionRequest
// }

type GetConversationsResponseModel struct {
	Data       GetConversationsData `json:"data"`
	Pagination Pagination           `json:"pagination"`
}

type GetConversationsData struct {
	Conversations []entity.ConversationList `json:"conversations"`
}

type ConversationResponseModel struct {
	MessageType string                   `json:"message_type"`
	Data        ConversationResponseData `json:"data"`
}

type ConversationResponseData struct {
	SenderUUID       string `json:"sender_uuid"`
	ConversationUUID string `json:"conversation_uuid"`
	SendMessageResponseData
	MessageDeletionConfirmation
	ReactionResponseData
	ErrorResponseData
}

type ErrorResponseData struct {
	ErrorMessage string `json:"error_msg"`
}
