package boundary

import "github.com/maxyong7/chat-messaging-app/internal/entity"

type InboxResponseModel struct {
	Data       InboxData  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type InboxData struct {
	Conversations []entity.Conversations `json:"conversations"`
}
