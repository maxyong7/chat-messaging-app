package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// LoginUseCase -.
type LoginUseCase struct {
	repo UserRepo
}

// New -.
func NewAuth(r UserRepo) *LoginUseCase {
	return &LoginUseCase{
		repo: r,
	}
}

// VerifyCredentials -.
func (uc *LoginUseCase) VerifyCredentials(ctx context.Context, v entity.UserCredentials) (string, bool, error) {
	userCredentialsDTO := entity.UserCredentialsDTO{
		Username: v.Username,
		Password: v.Password,
		Email:    v.Email,
	}
	userInfo, err := uc.repo.GetUserCredentials(ctx, userCredentialsDTO)
	if err != nil {
		return "", false, fmt.Errorf("LoginUseCase - VerifyCredentials - s.repo.GetUserInfo: %w", err)
	}

	if userInfo == nil {
		return "", false, entity.ErrUserNotFound
	}

	if userInfo.Password == v.Password {
		return userInfo.UserUuid, true, nil
	}
	return "", false, nil
}

// RegisterUser -.
func (uc *LoginUseCase) RegisterUser(ctx context.Context, userRegistration entity.UserRegistration) error {
	userRegistrationDTO := entity.UserRegistrationDTO{
		Username:  userRegistration.Username,
		Password:  userRegistration.Password,
		Email:     userRegistration.Email,
		FirstName: userRegistration.FirstName,
		LastName:  userRegistration.LastName,
		Avatar:    userRegistration.Avatar,
	}
	exist, err := uc.repo.CheckUserExist(ctx, userRegistrationDTO)
	if err != nil {
		return err
	}

	if exist {
		return entity.ErrUserAlreadyExists
	}

	err = uc.repo.StoreUserInfo(context.Background(), userRegistrationDTO)
	if err != nil {
		return fmt.Errorf("LoginUseCase - VerifyCredentials - s.repo.Store: %w", err)
	}
	return nil
}
