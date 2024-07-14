package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// ContactsUseCase -.
type ContactsUseCase struct {
	repo         ContactsRepo
	userInfoRepo UserRepo
}

// New -.
func NewContacts(r ContactsRepo, userInfoRepo UserRepo) *ContactsUseCase {
	return &ContactsUseCase{
		repo:         r,
		userInfoRepo: userInfoRepo,
	}
}

// VerifyCredentials -.
func (uc *ContactsUseCase) AddContacts(ctx context.Context, contactUserName string, userUuid string) error {
	// Check username exists
	contactUserUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, contactUserName)
	if err != nil {
		return fmt.Errorf("AddContacts - GetUserUUIDByUsername: %w", err)
	}

	if contactUserUUID == nil {
		return entity.ErrUserNameNotFound
	}

	// Check if already in contacts
	exist, err := uc.repo.CheckContactExist(ctx, userUuid, *contactUserUUID)
	if err != nil {
		return err
	}

	if exist {
		return entity.ErrUserAlreadyExists
	}

	//Store contacts
	contactsDTO := entity.ContactsDTO{
		UserUUID:         userUuid,
		ContactUserUUID:  *contactUserUUID,
		ConversationUUID: uuid.New().String(),
	}

	err = uc.repo.StoreContacts(ctx, contactsDTO)
	if err != nil {
		return fmt.Errorf("AddContacts - uc.repo.StoreContacts: %w", err)
	}
	return nil
}

// // RegisterUser -.
// func (uc *ContactsUseCase) RegisterUser(ctx context.Context, userRegistration entity.UserRegistration) error {
// 	exist, err := uc.repo.CheckUserExist(ctx, userRegistration)
// 	if err != nil {
// 		return err
// 	}

// 	if exist {
// 		return entity.ErrUserAlreadyExists
// 	}

// 	err = uc.repo.StoreUserInfo(context.Background(), userRegistration)
// 	if err != nil {
// 		return fmt.Errorf("LoginUseCase - VerifyCredentials - s.repo.Store: %w", err)
// 	}
// 	return nil
// }
