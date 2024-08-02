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
	userInfoDTO, err := uc.repo.GetUserProfile(ctx, userUUID)
	if err != nil {
		return userInfoDTO.ToUserInfo(), fmt.Errorf("UserProfileUseCase - GetUserInfo - GetUserInfo: %w", err)
	}
	if userInfoDTO == nil {
		return userInfoDTO.ToUserInfo(), entity.ErrUserNotFound
	}
	return userInfoDTO.ToUserInfo(), nil
}

func (uc *UserProfileUseCase) UpdateUserProfile(ctx context.Context, userInfo entity.UserProfile) error {
	userInfoDTO := entity.UserProfileDTO(userInfo)
	err := uc.repo.UpdateUserProfile(ctx, userInfoDTO)
	if err != nil {
		return fmt.Errorf("UserProfileUseCase - UpdateUserProfile - UpdateUserProfile: %w", err)
	}
	return nil
}
