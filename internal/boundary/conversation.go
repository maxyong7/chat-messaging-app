package boundary

import "encoding/json"

type ConversationRequestModel struct {
	MessageType string          `json:"message_type" binding:"required"`
	Data        json.RawMessage `json:"data" binding:"required"`
}

// type ConversationRequestData struct {
// 	SendMessageRequest
// 	DeleteMessageRequest
// 	AddReactionRequest
// }

type ConversationResponseModel struct {
	MessageType string                   `json:"message_type"`
	Data        ConversationResponseData `json:"data"`
}

type ConversationResponseData struct {
	SenderUUID       string `json:"sender_uuid"`
	ConversationUUID string `json:"conversation_uuid"`
	SendMessageResponseData
	DeleteMessageResponseData
	ReactionResponseData
	ErrorResponseData
}

type ErrorResponseData struct {
	ErrorMessage string `json:"error_msg"`
}
