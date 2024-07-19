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

func (uc *UserProfileUseCase) GetUserInfo(ctx context.Context, userUUID string) (*entity.UserInfoDTO, error) {
	userInfo, err := uc.repo.GetUserInfo(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("UserProfileUseCase - GetUserInfo - GetUserInfo: %w", err)
	}
	return userInfo, nil
}

func (uc *UserProfileUseCase) UpdateUserProfile(ctx context.Context, userInfoDTO entity.UserInfoDTO) error {
	err := uc.repo.UpdateUserInfo(ctx, userInfoDTO)
	if err != nil {
		return fmt.Errorf("UserProfileUseCase - UpdateUserProfile - UpdateUserProfile: %w", err)
	}
	return nil
}
