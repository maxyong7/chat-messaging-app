package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// UserProfileUseCase -.
type UserProfileUseCase struct {
	repo UserRepo
}

// New -.
func NewUserProfile(r UserRepo) *UserProfileUseCase {
	return &UserProfileUseCase{
		repo: r,
	}
}

func (uc *UserProfileUseCase) GetUserInfo(ctx context.Context, userUUID string) (entity.UserInfo, error) {
	userInfoDTO, err := uc.repo.GetUserInfo(ctx, userUUID)
	if err != nil {
		return userInfoDTO.ToUserInfo(), fmt.Errorf("UserProfileUseCase - GetUserInfo - GetUserInfo: %w", err)
	}
	if userInfoDTO == nil {
		return userInfoDTO.ToUserInfo(), entity.ErrUserNotFound
	}
	return userInfoDTO.ToUserInfo(), nil
}

func (uc *UserProfileUseCase) UpdateUserProfile(ctx context.Context, userInfo entity.UserInfo) error {
	userInfoDTO := entity.UserInfoDTO(userInfo)
	err := uc.repo.UpdateUserInfo(ctx, userInfoDTO)
	if err != nil {
		return fmt.Errorf("UserProfileUseCase - UpdateUserProfile - UpdateUserProfile: %w", err)
	}
	return nil
}
