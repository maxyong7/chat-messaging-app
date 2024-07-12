package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// LoginUseCase -.
type LoginUseCase struct {
	repo   UserRepo
	webAPI TranslationWebAPI
}

// New -.
func NewAuth(r UserRepo, w TranslationWebAPI) *LoginUseCase {
	return &LoginUseCase{
		repo:   r,
		webAPI: w,
	}
}

// VerifyCredentials -.
func (uc *LoginUseCase) VerifyCredentials(ctx context.Context, v entity.UserCredentials) (bool, error) {
	userInfo, err := uc.repo.GetUserInfo(ctx, v)
	if err != nil {
		return false, fmt.Errorf("LoginUseCase - VerifyCredentials - s.repo.GetUserInfo: %w", err)
	}

	if userInfo.Password == v.Password {
		return true, nil
	}
	return false, nil
}

// RegisterUser -.
func (uc *LoginUseCase) RegisterUser(ctx context.Context, userRegistration entity.UserRegistration) error {
	exist, err := uc.repo.CheckUserExist(ctx, userRegistration)
	if err != nil {
		return err
	}

	if exist {
		return entity.ErrUserAlreadyExists
	}

	err = uc.repo.StoreUserInfo(context.Background(), userRegistration)
	if err != nil {
		return fmt.Errorf("LoginUseCase - VerifyCredentials - s.repo.Store: %w", err)
	}
	return nil
}
