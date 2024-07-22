package boundary

import (
	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type MessageResponseModel struct {
	Data       MessageData `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type MessageData struct {
	Messages []entity.Message `json:"messages"`
}
