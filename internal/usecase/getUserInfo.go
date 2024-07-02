package usecase

import (
	"context"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// AuthenticateUserUseCase -.
type AuthenticateUserUseCase struct {
	repo   TranslationRepo
	webAPI TranslationWebAPI
}

// New -.
func NewAuth(r TranslationRepo, w TranslationWebAPI) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{
		repo:   r,
		webAPI: w,
	}
}

// History - getting translate history from store.
func (uc *AuthenticateUserUseCase) History(ctx context.Context) ([]entity.Translation, error) {
	translations, err := uc.repo.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("AuthenticateUserUseCase - History - s.repo.GetHistory: %w", err)
	}

	return translations, nil
}

// Translate -.
func (uc *AuthenticateUserUseCase) Translate(ctx context.Context, t entity.Translation) (entity.Translation, error) {
	translation, err := uc.webAPI.Translate(t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("AuthenticateUserUseCase - Translate - s.webAPI.Translate: %w", err)
	}

	err = uc.repo.Store(context.Background(), translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("AuthenticateUserUseCase - Translate - s.repo.Store: %w", err)
	}

	return translation, nil
}
