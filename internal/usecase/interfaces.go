// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

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
		// Store(context.Context, entity.Translation) error
		// GetHistory(context.Context) ([]entity.Translation, error)
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
		UpdateUserInfo(ctx context.Context, userInfo entity.UserInfo) error
	}

	Conversation interface {
		// ServeWs(*gin.Context, *Hub, string)
		// ServeWsWithRW(http.ResponseWriter, *http.Request, *Hub)
		StoreConversationAndMessage(ctx context.Context, conv entity.Conversation) error
	}

	ConversationRepo interface {
		GetConversations(context.Context, entity.RequestParams) ([]entity.Conversations, error)
		InsertConversationAndMessage(ctx context.Context, convDTO entity.ConversationDTO) error
	}

	Inbox interface {
		GetInbox(context.Context, entity.RequestParams) ([]entity.Conversations, error)
	}

	Contact interface {
		GetContacts(ctx context.Context, userUuid string) ([]entity.Contacts, error)
		AddContact(ctx context.Context, contactUserName string, userUuid string) error
		RemoveContact(ctx context.Context, contactUserName string, userUuid string) error
		UpdateBlockContact(ctx context.Context, contactUserName string, userUuid string, block bool) error
	}

	ContactsRepo interface {
		GetContactsByUserUUID(ctx context.Context, userUuid string) ([]entity.Contacts, error)
		CheckContactExist(ctx context.Context, userUuid string, contactUserUuid string) (bool, error)
		StoreContacts(context.Context, entity.ContactsDTO) error
		UpdateRemoved(context.Context, entity.ContactsDTO) error
		UpdateBlocked(context.Context, entity.ContactsDTO) error
	}

	MessageRepo interface {
		GetMessages(ctx context.Context, reqParam entity.RequestParams, conversationUUID string) ([]entity.GetMessageDTO, error)
		ValidateMessageSentByUser(ctx context.Context, msg entity.MessageDTO) (bool, error)
		DeleteMessage(ctx context.Context, msg entity.MessageDTO) error
	}

	Message interface {
		GetMessagesFromConversation(ctx context.Context, reqParam entity.RequestParams, conversationUUID string) ([]entity.GetMessageDTO, error)
		DeleteMessage(ctx context.Context, msg entity.Message) error
		ValidateMessageSentByUser(ctx context.Context, msg entity.Message) (bool, error)
	}

	ReactionRepo interface {
		StoreReaction(ctx context.Context, srDTO entity.StoreReactionDTO) error
		GetReactions(ctx context.Context, messageUUID string) ([]entity.GetReaction, error)
		RemoveReaction(ctx context.Context, rr entity.RemoveReactionDTO) error
	}

	Reaction interface {
		StoreReaction(ctx context.Context, reaction entity.Reaction) error
		RemoveReaction(ctx context.Context, reaction entity.Reaction) error
	}

	GroupChatRepo interface {
		CreateGroupChat(ctx context.Context, groupChat entity.GroupChatRequest) error
		AddParticipants(ctx context.Context, groupChat entity.GroupChatRequest) error
		RemoveParticipants(ctx context.Context, groupChat entity.GroupChatRequest) error
		UpdateGroupTitle(ctx context.Context, groupChat entity.GroupChatRequest) error
		ValidateUserInGroupChat(ctx context.Context, conversationUUID string, userUUID string) (bool, error)
	}

	GroupChat interface {
		CreateGroupChat(ctx context.Context, groupChatReq entity.GroupChatRequest) error
		AddParticipant(ctx context.Context, groupChatReq entity.GroupChatRequest) error
		RemoveParticipant(ctx context.Context, groupChatReq entity.GroupChatRequest) error
		UpdateGroupTitle(ctx context.Context, groupChatReq entity.GroupChatRequest) error
	}

	UserProfile interface {
		GetUserInfo(ctx context.Context, userUUID string) (entity.UserInfo, error)
		UpdateUserProfile(ctx context.Context, userInfoDTO entity.UserInfo) error
	}
)
