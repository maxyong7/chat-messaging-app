package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

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

func (uc *ConversationUseCase) GetConversations(ctx context.Context, reqParam entity.RequestParams) ([]entity.Conversations, error) {
	reqParamDTO := entity.RequestParamsDTO(reqParam)
	conversations, err := uc.repo.GetConversations(ctx, reqParamDTO)
	if err != nil {
		return nil, fmt.Errorf("ConversationUseCase - GetConversation - s.repo.GetConversations: %w", err)
	}
	if len(conversations) == 0 {
		return nil, nil
	}

	return conversations, nil
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
