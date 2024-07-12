package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// InboxUseCase -.
type InboxUseCase struct {
	repo ConversationRepo
}

// New -.
func NewInbox(r ConversationRepo) *InboxUseCase {
	return &InboxUseCase{
		repo: r,
	}
}

// VerifyCredentials -.
func (uc *InboxUseCase) GetInbox(ctx context.Context, reqParam entity.RequestParams) (entity.InboxResponse, error) {
	conversations, err := uc.repo.GetConversations(ctx, reqParam)
	if err != nil {
		return entity.InboxResponse{}, fmt.Errorf("InboxUseCase - GetInbox - s.repo.GetConversations: %w", err)
	}
	if len(conversations) == 0 {
		return entity.InboxResponse{}, nil
	}

	encodedCursor := encodeCursor(conversations[len(conversations)-1].LastMessageCreatedAt)
	return entity.InboxResponse{
		Data: entity.Data{
			Conversations: conversations,
		},
		Pagination: entity.Pagination{
			Cursor: encodedCursor,
			Limit:  reqParam.Limit,
		},
	}, nil
}
