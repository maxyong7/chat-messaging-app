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
	// Convert request parameter entity object into reqParamDTO
	reqParamDTO := entity.RequestParamsDTO(reqParam)

	// Use reqParamDTO to get a conversation list from conversation data repository
	// It queries 'contacts' (direct message) table and 'participants' (group message) table
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
	// Convert conversation entity object into convDTO
	convDTO := entity.ConversationDTO{
		SenderUUID:       conv.SenderUUID,
		ConversationUUID: conv.ConversationUUID,
		MessageUUID:      conv.MessageUUID,
		Content:          conv.Content,
		CreatedAt:        conv.CreatedAt,
	}

	// Insert message into 'messages' table to store the entire conversation history
	// and Upsert 'conversations' table with the most recent message
	// using conversation data repository
	err := uc.repo.InsertConversationAndMessage(ctx, convDTO)
	if err != nil {
		return fmt.Errorf("ConversationUseCase - StoreConversation - uc.repo.InsertConversationAndMessage: %w", err)
	}
	return nil
}
