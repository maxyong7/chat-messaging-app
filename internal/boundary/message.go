package boundary

import (
	"time"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type SendMessageRequest struct {
	Content string `json:"content"`
}

type DeleteMessageRequest struct {
	MessageUUID string `json:"message_uuid"`
}

type SendMessageResponseData struct {
	SenderFirstName string    `json:"sender_first_name"`
	SenderLastName  string    `json:"sender_last_name"`
	Content         string    `json:"content"`
	MessageUUID     string    `json:"message_uuid"`
	CreatedAt       time.Time `json:"created_at"`
}
type DeleteMessageResponseData struct {
	MessageUUID string `json:"message_uuid"`
}

type GetMessageResponseModel struct {
	Data       GetMessageResponseData `json:"data"`
	Pagination Pagination             `json:"pagination"`
}

type GetMessageResponseData struct {
	Messages []entity.GetMessageDTO `json:"messages"`
}
