package boundary

import (
	"time"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type ChatInterface struct {
	Content string `json:"content"`
}

type DeleteMessageRequest struct {
	MessageUUID string `json:"message_uuid"`
}

type SendMessageResponseData struct {
	SenderFirstName string    `json:"sender_first_name"`
	SenderLastName  string    `json:"sender_last_name"`
	SenderAvatar    string    `json:"sender_avatar"`
	Content         string    `json:"content"`
	MessageUUID     string    `json:"message_uuid"`
	CreatedAt       time.Time `json:"created_at"`
}
type MessageDeletionConfirmation struct {
	MessageUUID string `json:"message_uuid"`
}

type ConversationScreen struct {
	Data       ConversationData `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

type ConversationData struct {
	Messages []entity.GetMessageDTO `json:"messages"`
}

type MessageStatusIndicator struct {
	SeenStatus []entity.GetSeenStatusDTO `json:"seen_status"`
}

type MessageSearchScreen struct {
	Data MessageSearchData `json:"data"`
}

type MessageSearchData struct {
	Messages []entity.SearchMessageDTO `json:"messages"`
}
