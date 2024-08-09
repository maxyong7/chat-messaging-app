package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type GroupChatUseCase struct {
	repo         GroupChatRepo
	userInfoRepo UserRepo
}

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
	// Convert groupChat entity object into groupChatDTO
	groupChatDTO := toGroupChatDTO(groupChat)

	// Check if user is in groupchat from group chat data repository
	exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatDTO.ConversationUUID, groupChatDTO.UserUUID)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - AddParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
	}
	// If user not in group chat, dont allow user to add participant
	if !exist {
		return entity.ErrUserNotInGroupChat
	}

	for i := range groupChatDTO.Participants {
		// Check if participant is in groupchat from group chat data repository
		exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatDTO.ConversationUUID, groupChatDTO.Participants[i].ParticipantUUID)
		if err != nil {
			return fmt.Errorf("GroupChatUseCase - AddParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
		}
		if exist {
			return entity.ErrParticipantAlrdInGroupChat
		}
	}

	// Add participants in group chat data repository using groupChatDTO
	err = uc.repo.AddParticipants(ctx, groupChatDTO)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - AddParticipant - uc.repo.AddParticipants: %w", err)
	}
	return nil
}

func (uc *GroupChatUseCase) RemoveParticipant(ctx context.Context, groupChat entity.GroupChat) error {
	// Convert groupChat entity object into groupChatDTO
	groupChatDTO := toGroupChatDTO(groupChat)

	// Check if user is in groupchat from group chat data repository
	exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatDTO.ConversationUUID, groupChatDTO.UserUUID)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - RemoveParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
	}
	// If user not in group chat, dont allow user to remove participant
	if !exist {
		return entity.ErrUserNotInGroupChat
	}

	for i := range groupChatDTO.Participants {
		// Check if participant is in groupchat from group chat data repository
		exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatDTO.ConversationUUID, groupChatDTO.Participants[i].ParticipantUUID)
		if err != nil {
			return fmt.Errorf("GroupChatUseCase - RemoveParticipant - uc.repo.ValidateUserInGroupChat: %w", err)
		}
		if !exist {
			return entity.ErrParticipantNotInGroupChat
		}
	}

	// Remove participants in group chat data repository using groupChatDTO
	err = uc.repo.RemoveParticipants(ctx, groupChatDTO)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - RemoveParticipant - uc.repo.RemoveParticipants: %w", err)
	}
	return nil
}

func (uc *GroupChatUseCase) UpdateGroupTitle(ctx context.Context, groupChat entity.GroupChat) error {
	// Convert groupChat entity object into groupChatDTO
	groupChatDTO := toGroupChatDTO(groupChat)

	// Check if user is in groupchat from group chat data repository
	exist, err := uc.repo.ValidateUserInGroupChat(ctx, groupChatDTO.ConversationUUID, groupChatDTO.UserUUID)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - UpdateGroupTitle - uc.repo.ValidateUserInGroupChat: %w", err)
	}
	// If user not in group chat, dont allow user to update group title
	if !exist {
		return entity.ErrUserNotInGroupChat
	}

	// Update group title in group chat data repository using groupChatDTO
	err = uc.repo.UpdateGroupTitle(ctx, groupChatDTO)
	if err != nil {
		return fmt.Errorf("GroupChatUseCase - UpdateGroupTitle - uc.repo.UpdateGroupTitles: %w", err)
	}
	return nil
}

// Convert group chat entity object to group chat DTO
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
