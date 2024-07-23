package boundary

type ConversationRequestModel struct {
	MessageType string                  `json:"message_type" binding:"required"`
	Data        ConversationRequestData `json:"data" binding:"required"`
}

type ConversationRequestData struct {
	SendMessageRequest   SendMessageRequest
	DeleteMessageRequest DeleteMessageRequest
	AddReactionRequest   AddReactionRequest
}

type ConversationResponseModel struct {
	MessageType string                   `json:"message_type"`
	Data        ConversationResponseData `json:"data"`
}

type ConversationResponseData struct {
	SenderUUID                string `json:"sender_uuid"`
	ConversationUUID          string `json:"conversation_uuid"`
	SendMessageResponseData   SendMessageResponseData
	DeleteMessageResponseData DeleteMessageResponseData
	AddReactionResponseData   ReactionResponseData
	ErrorResponseData         ErrorResponseData
}

type ErrorResponseData struct {
	ErrorMessage string `json:"error_msg"`
}
