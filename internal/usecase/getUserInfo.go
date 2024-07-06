package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// AuthenticateUserUseCase -.
type AuthenticateUserUseCase struct {
	repo   VerificationRepo
	webAPI TranslationWebAPI
}

// New -.
func NewAuth(r VerificationRepo, w TranslationWebAPI) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{
		repo:   r,
		webAPI: w,
	}
}

// // History - getting translate history from store.
// func (uc *AuthenticateUserUseCase) History(ctx context.Context) ([]entity.Translation, error) {
// 	translations, err := uc.repo.GetHistory(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("AuthenticateUserUseCase - History - s.repo.GetHistory: %w", err)
// 	}

// 	return translations, nil
// }

// VerifyCredentials -.
func (uc *AuthenticateUserUseCase) VerifyCredentials(ctx context.Context, v entity.Verification) (bool, error) {
	userInfo, err := uc.repo.GetUserInfo(ctx, v)
	if err != nil {
		return false, fmt.Errorf("AuthenticateUserUseCase - VerifyCredentials - s.repo.GetUserInfo: %w", err)
	}

	if userInfo.Password == v.Password {
		return true, nil
	}

	// err = uc.repo.Store(context.Background(), translation)
	// if err != nil {
	// 	return false, fmt.Errorf("AuthenticateUserUseCase - VerifyCredentials - s.repo.Store: %w", err)
	// }

	return false, nil
}
