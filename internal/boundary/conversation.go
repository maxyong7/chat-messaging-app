package boundary

import "github.com/maxyong7/chat-messaging-app/internal/entity"

type ConversationRequestModel struct {
	MessageType string                  `json:"message_type" binding:"required"`
	Data        ConversationRequestData `json:"data" binding:"required"`
}

type ConversationRequestData struct {
	MessageRequestData
}

func (r ConversationRequestData) ToMessage(conversationUUID string) entity.ConversationMessage {
	return entity.ConversationMessage{
		ConversationUUID: conversationUUID,
		Content:          r.Content,
	}
}

type ConversationResponseModel struct {
	MessageType string                   `json:"message_type"`
	Data        ConversationResponseData `json:"data"`
}

type ConversationResponseData struct {
	Conversation entity.Conversation
}
