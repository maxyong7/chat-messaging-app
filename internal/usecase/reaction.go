package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type ReactionUseCase struct {
	reactionRepo ReactionRepo
}

func NewReaction(r ReactionRepo) *ReactionUseCase {
	return &ReactionUseCase{
		reactionRepo: r,
	}
}

func (uc *ReactionUseCase) StoreReaction(ctx context.Context, reaction entity.Reaction) error {
	// Convert reaction entity object into storeReactionDTO
	storeReactionDTO := entity.StoreReactionDTO{
		MessageUUID:  reaction.MessageUUID,
		SenderUUID:   reaction.SenderUUID,
		ReactionType: reaction.ReactionType,
	}

	// Store reaction into 'reaction' table using reaction data repository
	err := uc.reactionRepo.StoreReaction(ctx, storeReactionDTO)
	if err != nil {
		return fmt.Errorf("ReactionUseCase - StoreReaction - uc.reactionRepo.StoreReaction: %w", err)
	}
	return nil
}

func (uc *ReactionUseCase) RemoveReaction(ctx context.Context, reaction entity.Reaction) error {
	// Convert reaction entity object into removeReactionDTO
	removeReactionDTO := entity.RemoveReactionDTO{
		MessageUUID: reaction.MessageUUID,
		SenderUUID:  reaction.SenderUUID,
	}

	// Remove reaction from 'reaction' table using reaction data repository
	err := uc.reactionRepo.RemoveReaction(ctx, removeReactionDTO)
	if err != nil {
		return fmt.Errorf("ReactionUseCase - RemoveReaction - uc.reactionRepo.RemoveReaction: %w", err)
	}
	return nil
}
