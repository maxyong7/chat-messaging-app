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
func (uc *InboxUseCase) GetInbox(ctx context.Context, reqParam entity.RequestParams) ([]entity.Conversations, error) {
	conversations, err := uc.repo.GetConversations(ctx, reqParam)
	if err != nil {
		return nil, fmt.Errorf("InboxUseCase - GetInbox - s.repo.GetConversations: %w", err)
	}
	if len(conversations) == 0 {
		return nil, nil
	}

	return conversations, nil
}
