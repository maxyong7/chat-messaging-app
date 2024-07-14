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

func (uc *ContactsUseCase) GetContacts(ctx context.Context, userUuid string) ([]entity.Contacts, error) {
	contacts, err := uc.repo.GetContactsByUserUUID(ctx, userUuid)
	if err != nil {
		return nil, fmt.Errorf("ContactsUseCase - GetContacts - GetContactsByUserUUID: %w", err)
	}
	return contacts, nil
}

func (uc *ContactsUseCase) AddContact(ctx context.Context, contactUserName string, userUuid string) error {
	// Check username exists
	contactUserUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, contactUserName)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - AddContacts - GetUserUUIDByUsername: %w", err)
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
		err := uc.repo.UpdateRemoved(ctx, entity.ContactsDTO{
			UserUUID:        userUuid,
			ContactUserUUID: *contactUserUUID,
			Removed:         false,
		})

		if err != nil {
			return fmt.Errorf("ContactsUseCase - AddContacts - uc.repo.UpdateRemoved: %w", err)
		}
		return nil
	}

	//Store contacts
	contactsDTO := entity.ContactsDTO{
		UserUUID:         userUuid,
		ContactUserUUID:  *contactUserUUID,
		ConversationUUID: uuid.New().String(),
	}
	err = uc.repo.StoreContacts(ctx, contactsDTO)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - AddContacts - uc.repo.StoreContacts: %w", err)
	}
	return nil
}

func (uc *ContactsUseCase) RemoveContact(ctx context.Context, contactUserName string, userUuid string) error {
	// Check username exists
	contactUserUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, contactUserName)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - AddContacts - GetUserUUIDByUsername: %w", err)
	}

	if contactUserUUID == nil {
		return entity.ErrUserNameNotFound
	}

	// Check if already in contacts
	exist, err := uc.repo.CheckContactExist(ctx, userUuid, *contactUserUUID)
	if err != nil {
		return err
	}

	if !exist {
		return entity.ErrContactDoesNotExists
	}

	//Store contacts
	contactsDTO := entity.ContactsDTO{
		UserUUID:        userUuid,
		ContactUserUUID: *contactUserUUID,
		Removed:         true,
	}

	err = uc.repo.UpdateRemoved(ctx, contactsDTO)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - AddContacts - uc.repo.StoreContacts: %w", err)
	}
	return nil
}
