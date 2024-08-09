package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type UserProfileUseCase struct {
	repo UserRepo
}

// New -.
func NewUserProfile(r UserRepo) *UserProfileUseCase {
	return &UserProfileUseCase{
		repo: r,
	}
}

func (uc *UserProfileUseCase) GetUserProfile(ctx context.Context, userUUID string) (entity.UserProfile, error) {
	// Get user profile from user profile data repository
	userProfileDTO, err := uc.repo.GetUserProfile(ctx, userUUID)

	//Convert userProfileDTO into user info entity object
	userProfileEntity := userProfileDTO.ToUserInfo()
	if err != nil {
		return userProfileEntity, fmt.Errorf("UserProfileUseCase - GetUserInfo - GetUserInfo: %w", err)
	}
	if userProfileDTO == nil {
		return userProfileEntity, entity.ErrUserNotFound
	}
	return userProfileEntity, nil
}

func (uc *UserProfileUseCase) UpdateUserProfile(ctx context.Context, userProfile entity.UserProfile) error {
	// Convert user profile entity object into userProfileDTO
	userProfileDTO := entity.UserProfileDTO(userProfile)

	// Update user profile in user profile data repository
	err := uc.repo.UpdateUserProfile(ctx, userProfileDTO)
	if err != nil {
		return fmt.Errorf("UserProfileUseCase - UpdateUserProfile - UpdateUserProfile: %w", err)
	}
	return nil
}
