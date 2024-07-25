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

func (uc *GroupChatUseCase) CreateGroupChat(ctx context.Context, groupChat entity.GroupChat) error {
	err := uc.repo.CreateGroupChat(ctx, toGroupChatDTO(groupChat))
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - CreateGroupChat - CreateGroupChat: %w", err)
	}
	return nil
}

func (uc *GroupChatUseCase) AddParticipant(ctx context.Context, groupChat entity.GroupChat) error {
	// Check if user is in groupchat before allowing to add participant
	groupChatDTO := toGroupChatDTO(groupChat)
	exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatDTO.ConversationUUID, groupChatDTO.UserUUID)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - AddParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
	}
	if !exist {
		return entity.ErrUserNotInGroupChat
	}

	for i := range groupChatDTO.Participants {
		// Check if participant is in groupchat
		exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatDTO.ConversationUUID, groupChatDTO.Participants[i].ParticipantUUID)
		if err != nil {
			return fmt.Errorf("GroupChatUseCase - AddParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
		}
		if exist {
			return entity.ErrParticipantAlrdInGroupChat
		}
	}

	// // Check username exists [Needed if front end can only pass username, not participant's uuid]
	// for i, participant := range groupChat.Participants {
	// 	participantUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, participant.Username)
	// 	if err != nil {
	// 		return fmt.Errorf("GroupChatUseCase - AddParticipants - GetUserUUIDByUsername: %w", err)
	// 	}

	// 	if participantUUID == nil {
	// 		return entity.ErrUserNameNotFound
	// 	}

	// 	groupChat.Participants[i].ParticipantUUID = *participantUUID
	// }

	err = uc.repo.AddParticipants(ctx, groupChatDTO)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - AddParticipant - uc.repo.AddParticipants: %w", err)
	}
	return nil
}

func (uc *GroupChatUseCase) RemoveParticipant(ctx context.Context, groupChat entity.GroupChat) error {
	groupChatDTO := toGroupChatDTO(groupChat)
	// Check if user is in groupchat before allowing to remove participant
	exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatDTO.ConversationUUID, groupChatDTO.UserUUID)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - RemoveParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
	}
	if !exist {
		return entity.ErrUserNotInGroupChat
	}

	for i := range groupChatDTO.Participants {
		// Check if participant is in groupchat
		exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatDTO.ConversationUUID, groupChatDTO.Participants[i].ParticipantUUID)
		if err != nil {
			return fmt.Errorf("GroupChatUseCase - RemoveParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
		}
		if !exist {
			return entity.ErrParticipantNotInGroupChat
		}
	}

	err = uc.repo.RemoveParticipants(ctx, groupChatDTO)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - RemoveParticipant - uc.repo.RemoveParticipants: %w", err)
	}
	return nil
}

func (uc *GroupChatUseCase) UpdateGroupTitle(ctx context.Context, groupChat entity.GroupChat) error {
	groupChatDTO := toGroupChatDTO(groupChat)
	// Check if user is in groupchat before allowing to remove participant
	exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatDTO.ConversationUUID, groupChatDTO.UserUUID)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - UpdateGroupTitle - uc.repo.ValidateUserInGroupChat: %w", err)
	}
	if !exist {
		return entity.ErrUserNotInGroupChat
	}

	err = uc.repo.UpdateGroupTitle(ctx, groupChatDTO)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - UpdateGroupTitle - uc.repo.UpdateGroupTitles: %w", err)
	}
	return nil
}

func toGroupChatDTO(gc entity.GroupChat) entity.GroupChatDTO {
	participantsDTO := []entity.ParticipantDTO{}
	for _, p := range gc.Participants {
		participantDTO := entity.ParticipantDTO{
			Username:        p.Username,
			ParticipantUUID: p.ParticipantUUID,
		}
		participantsDTO = append(participantsDTO, participantDTO)

	}
	return entity.GroupChatDTO{
		UserUUID:         gc.UserUUID,
		Title:            gc.Title,
		ConversationUUID: gc.ConversationUUID,
		Participants:     participantsDTO,
	}
}
