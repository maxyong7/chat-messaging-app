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

func (uc *MessageUseCase) GetMessagesFromConversation(ctx context.Context, reqParam entity.RequestParams, conversationUUID string) ([]entity.GetMessageDTO, error) {
	messages, err := uc.msgRepo.GetMessages(ctx, reqParam, conversationUUID)
	if err != nil {
		return nil, fmt.Errorf("MessageUseCase - GetMessages - uc.msgRepo.GetMessages: %w", err)
	}
	if len(messages) == 0 {
		return nil, nil
	}

	for i, msg := range messages {
		reactions, err := uc.reactionRepo.GetReactions(ctx, msg.MessageUUID)
		if err != nil {
			return nil, fmt.Errorf("MessageUseCase - GetMessages - uc.reactionRepo.GetReactions: %w", err)
		}
		messages[i].Reaction = reactions
	}

	return messages, nil
}

func (uc *MessageUseCase) DeleteMessage(ctx context.Context, conv entity.Conversation, msg entity.Message) error {
	err := uc.msgRepo.DeleteMessage(ctx, conv, msg)
	if err != nil {
		return fmt.Errorf("MessageUseCase - DeleteMessage - uc.msgRepo.DeleteMessage: %w", err)
	}
	return nil
}

func (uc *MessageUseCase) ValidateMessageSentByUser(ctx context.Context, conv entity.Conversation, msg entity.Message) (bool, error) {
	valid, err := uc.msgRepo.ValidateMessageSentByUser(ctx, conv, msg)
	if err != nil {
		return false, fmt.Errorf("MessageUseCase - DeleteMessage - uc.msgRepo.DeleteMessage: %w", err)
	}
	return valid, nil
}
