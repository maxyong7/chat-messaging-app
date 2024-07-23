package boundary

import (
	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type MessageRequestData struct {
	Content string `json:"content"`
}

type MessageResponseModel struct {
	Data       MessageResponseData `json:"data"`
	Pagination Pagination          `json:"pagination"`
}

type MessageResponseData struct {
	Messages []entity.Message `json:"messages"`
}
