package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// MessageUseCase -.
type MessageUseCase struct {
	msgRepo      MessageRepo
	reactionRepo ReactionRepo
}

// New -.
func NewMessage(m MessageRepo, r ReactionRepo) *MessageUseCase {
	return &MessageUseCase{
		msgRepo:      m,
		reactionRepo: r,
	}
}

// VerifyCredentials -.
func (uc *MessageUseCase) GetMessagesFromConversation(ctx context.Context, reqParam entity.RequestParams, conversationUUID string) (entity.MessageResponse, error) {
	messages, err := uc.msgRepo.GetMessages(ctx, reqParam, conversationUUID)
	if err != nil {
		return entity.MessageResponse{}, fmt.Errorf("MessageUseCase - GetMessages - uc.msgRepo.GetMessages: %w", err)
	}
	if len(messages) == 0 {
		return entity.MessageResponse{}, nil
	}

	for i, msg := range messages {
		reactions, err := uc.reactionRepo.GetReactions(ctx, msg.MessageUUID)
		if err != nil {
			return entity.MessageResponse{}, fmt.Errorf("MessageUseCase - GetMessages - uc.reactionRepo.GetReactions: %w", err)
		}
		messages[i].Reaction = reactions
	}

	encodedCursor := encodeCursor(&messages[len(messages)-1].CreatedAt)
	return entity.MessageResponse{
		Data: entity.MessageData{
			Messages: messages,
		},
		Pagination: entity.Pagination{
			Cursor: encodedCursor,
			Limit:  reqParam.Limit,
		},
	}, nil
}
