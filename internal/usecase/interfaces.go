// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Translation -.
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		History(context.Context) ([]entity.Translation, error)
	}

	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}

	// Verification -.
	Verification interface {
		VerifyCredentials(context.Context, entity.UserInfo) (bool, error)
		RegisterUser(context.Context, entity.UserInfo) error
	}

	// UserRepo -.
	UserRepo interface {
		GetUserInfo(context.Context, entity.UserInfo) (*entity.UserInfoDTO, error)
		StoreUserInfo(context.Context, entity.UserInfo) error
		CheckUserExist(context.Context, entity.UserInfo) (bool, error)
	}

	Conversation interface {
		ServeWs(*gin.Context, *Hub)
		ServeWsWithRW(http.ResponseWriter, *http.Request, *Hub)
	}

	ConversationRepo interface {
		GetConversations(context.Context, entity.RequestParams) ([]entity.Conversations, error)
	}

	Inbox interface {
		GetInbox(context.Context, entity.RequestParams) (entity.InboxResponse, error)
	}
)
