package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// ConversationUseCase -.
type ConversationUseCase struct {
	repo ConversationRepo
}

type ReactionData struct {
	ReactionType string `json:"reaction_type"`
}

// New -.
func NewConversation(r ConversationRepo) *ConversationUseCase {
	return &ConversationUseCase{
		repo: r,
	}
}

func (uc *ConversationUseCase) StoreConversationAndMessage(ctx context.Context, conv entity.Conversation) error {
	convDTO := entity.ConversationDTO{
		SenderUUID:       conv.SenderUUID,
		ConversationUUID: conv.ConversationUUID,
		MessageUUID:      conv.MessageUUID,
		Content:          conv.Content,
		CreatedAt:        conv.CreatedAt,
	}
	err := uc.repo.InsertConversationAndMessage(ctx, convDTO)
	if err != nil {
		return fmt.Errorf("ConversationUseCase - StoreConversation - uc.repo.InsertConversationAndMessage: %w", err)
	}
	return nil
}
