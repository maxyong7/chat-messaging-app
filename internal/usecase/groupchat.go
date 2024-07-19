package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// GroupChatUseCase -.
type GroupChatUseCase struct {
	repo         GroupChatRepo
	userInfoRepo UserRepo
}

// New -.
func NewGroupChat(r GroupChatRepo, userInfoRepo UserRepo) *GroupChatUseCase {
	return &GroupChatUseCase{
		repo:         r,
		userInfoRepo: userInfoRepo,
	}
}

func (uc *GroupChatUseCase) CreateGroupChat(ctx context.Context, groupChatReq entity.GroupChatRequest) error {
	err := uc.repo.CreateGroupChat(ctx, groupChatReq)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - CreateGroupChat - CreateGroupChat: %w", err)
	}
	return nil
}

func (uc *GroupChatUseCase) AddParticipant(ctx context.Context, groupChatReq entity.GroupChatRequest) error {
	// Check if user is in groupchat before allowing to add participant
	exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatReq.ConversationUUID, groupChatReq.UserUUID)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - AddParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
	}
	if !exist {
		return entity.ErrUserNotInGroupChat
	}

	for i := range groupChatReq.Participants {
		// Check if participant is in groupchat
		exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatReq.ConversationUUID, groupChatReq.Participants[i].ParticipantUUID)
		if err != nil {
			return fmt.Errorf("GroupChatUseCase - AddParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
		}
		if exist {
			return entity.ErrParticipantAlrdInGroupChat
		}
	}

	// // Check username exists [Needed if front end can only pass username, not participant's uuid]
	// for i, participant := range groupChatReq.Participants {
	// 	participantUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, participant.Username)
	// 	if err != nil {
	// 		return fmt.Errorf("GroupChatUseCase - AddParticipants - GetUserUUIDByUsername: %w", err)
	// 	}

	// 	if participantUUID == nil {
	// 		return entity.ErrUserNameNotFound
	// 	}

	// 	groupChatReq.Participants[i].ParticipantUUID = *participantUUID
	// }

	err = uc.repo.AddParticipants(ctx, groupChatReq)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - AddParticipant - uc.repo.AddParticipants: %w", err)
	}
	return nil
}

func (uc *GroupChatUseCase) RemoveParticipant(ctx context.Context, groupChatReq entity.GroupChatRequest) error {
	// Check if user is in groupchat before allowing to remove participant
	exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatReq.ConversationUUID, groupChatReq.UserUUID)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - RemoveParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
	}
	if !exist {
		return entity.ErrUserNotInGroupChat
	}

	for i := range groupChatReq.Participants {
		// Check if participant is in groupchat
		exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatReq.ConversationUUID, groupChatReq.Participants[i].ParticipantUUID)
		if err != nil {
			return fmt.Errorf("GroupChatUseCase - RemoveParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
		}
		if !exist {
			return entity.ErrParticipantNotInGroupChat
		}
	}

	err = uc.repo.RemoveParticipants(ctx, groupChatReq)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - RemoveParticipant - uc.repo.RemoveParticipants: %w", err)
	}
	return nil
}

func (uc *GroupChatUseCase) UpdateGroupTitle(ctx context.Context, groupChatReq entity.GroupChatRequest) error {
	// Check if user is in groupchat before allowing to remove participant
	exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatReq.ConversationUUID, groupChatReq.UserUUID)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - UpdateGroupTitle - uc.repo.ValidateUserInGroupChat: %w", err)
	}
	if !exist {
		return entity.ErrUserNotInGroupChat
	}

	err = uc.repo.UpdateGroupTitle(ctx, groupChatReq)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - UpdateGroupTitle - uc.repo.UpdateGroupTitles: %w", err)
	}
	return nil
}
