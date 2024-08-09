package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase struct {
	repo UserRepo
}

func NewAuth(r UserRepo) *LoginUseCase {
	return &LoginUseCase{
		repo: r,
	}
}

func (uc *LoginUseCase) VerifyCredentials(ctx context.Context, userCredentials entity.UserCredentials) (string, bool, error) {
	// Convert user credential entity object into userCredentialsDTO
	userCredentialsDTO := entity.UserCredentialsDTO{
		Username: userCredentials.Username,
		Password: userCredentials.Password,
		Email:    userCredentials.Email,
	}

	// Get user credentials from user data repository
	userInfo, err := uc.repo.GetUserCredentials(ctx, userCredentialsDTO)
	if err != nil {
		return "", false, fmt.Errorf("LoginUseCase - VerifyCredentials - s.repo.GetUserInfo: %w", err)
	}

	if userInfo == nil {
		return "", false, entity.ErrUserNotFound
	}

	if userInfo.Password == userCredentials.Password {
		return userInfo.UserUuid, true, nil
	}

	// hashedPassword, err := hashPassword(userCredentials.Password)
	// if err != nil {
	// 	return "", false, err
	// }
	// Verify if password matches data repository
	match := verifyPassword(userCredentials.Password, userInfo.Password)
	if match {
		return userInfo.UserUuid, true, nil
	}
	return "", false, entity.ErrIncorrectPassword
}

func (uc *LoginUseCase) RegisterUser(ctx context.Context, userRegistration entity.UserRegistration) error {
	hashedPassword, err := hashPassword(userRegistration.Password)
	if err != nil {
		return err
	}
	// Convert user registration entity object into userRegistrationDTO
	userRegistrationDTO := entity.UserRegistrationDTO{
		Username:  userRegistration.Username,
		Password:  hashedPassword,
		Email:     userRegistration.Email,
		FirstName: userRegistration.FirstName,
		LastName:  userRegistration.LastName,
		Avatar:    userRegistration.Avatar,
	}

	// Check if user already register from user data repository
	exist, err := uc.repo.CheckUserExist(ctx, userRegistrationDTO)
	if err != nil {
		return err
	}

	if exist {
		return entity.ErrUserAlreadyExists
	}

	// Register user in user data repository
	err = uc.repo.StoreUserInfo(context.Background(), userRegistrationDTO)
	if err != nil {
		return fmt.Errorf("LoginUseCase - VerifyCredentials - s.repo.Store: %w", err)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
