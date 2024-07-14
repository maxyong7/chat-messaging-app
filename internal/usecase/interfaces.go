// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

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
		VerifyCredentials(context.Context, entity.UserCredentials) (string, bool, error)
		RegisterUser(context.Context, entity.UserRegistration) error
	}

	// UserRepo -.
	UserRepo interface {
		GetUserCredentials(context.Context, entity.UserCredentials) (*entity.UserCredentialsDTO, error)
		StoreUserInfo(context.Context, entity.UserRegistration) error
		CheckUserExist(context.Context, entity.UserRegistration) (bool, error)
		GetUserInfo(context.Context, string) (*entity.UserInfoDTO, error)
		GetUserUUIDByUsername(context.Context, string) (*string, error)
	}

	Conversation interface {
		ServeWs(*gin.Context, *Hub, string)
		// ServeWsWithRW(http.ResponseWriter, *http.Request, *Hub)
	}

	ConversationRepo interface {
		GetConversations(context.Context, entity.RequestParams) ([]entity.Conversations, error)
		StoreConversation(Message) error
	}

	Inbox interface {
		GetInbox(context.Context, entity.RequestParams) (entity.InboxResponse, error)
	}

	Contact interface {
		AddContacts(ctx context.Context, contactUserName string, userUuid string) error
	}

	ContactsRepo interface {
		GetContactsByUserUUID(ctx context.Context, userUuid string) ([]entity.Contacts, error)
		CheckContactExist(ctx context.Context, userUuid string, contactUserUuid string) (bool, error)
		StoreContacts(context.Context, entity.ContactsDTO) error
	}
)
