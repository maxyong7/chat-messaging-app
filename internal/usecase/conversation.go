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

func NewConversation(r ConversationRepo) *ConversationUseCase {
	return &ConversationUseCase{
		repo: r,
	}
}

func (uc *ConversationUseCase) GetConversationList(ctx context.Context, reqParam entity.RequestParams) ([]entity.ConversationList, error) {
	// Convert request parameter entity object into DTO
	reqParamDTO := entity.RequestParamsDTO(reqParam)

	// Use converted DTO to get conversation list from conversation data repository
	conversations, err := uc.repo.GetConversationList(ctx, reqParamDTO)
	if err != nil {
		return nil, fmt.Errorf("ConversationUseCase - GetConversation - s.repo.GetConversationList: %w", err)
	}

	if len(conversations) == 0 {
		return nil, nil
	}

	return conversations, nil
}

func (uc *ConversationUseCase) StoreConversationAndMessage(ctx context.Context, conv entity.Conversation) error {
	// Convert conversation entity object into DTO
	convDTO := entity.ConversationDTO{
		SenderUUID:       conv.SenderUUID,
		ConversationUUID: conv.ConversationUUID,
		MessageUUID:      conv.MessageUUID,
		Content:          conv.Content,
		CreatedAt:        conv.CreatedAt,
	}

	// Insert conversation and message into conversation data repository
	err := uc.repo.InsertConversationAndMessage(ctx, convDTO)
	if err != nil {
		return fmt.Errorf("ConversationUseCase - StoreConversation - uc.repo.InsertConversationAndMessage: %w", err)
	}
	return nil
}
