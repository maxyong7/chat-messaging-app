package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

type ContactsUseCase struct {
	repo         ContactsRepo
	userInfoRepo UserRepo
}

func NewContacts(r ContactsRepo, userInfoRepo UserRepo) *ContactsUseCase {
	return &ContactsUseCase{
		repo:         r,
		userInfoRepo: userInfoRepo,
	}
}

func (uc *ContactsUseCase) GetContacts(ctx context.Context, userUuid string) ([]entity.Contacts, error) {
	// Get contacts from data repository
	contacts, err := uc.repo.GetContactsByUserUUID(ctx, userUuid)
	if err != nil {
		return nil, fmt.Errorf("ContactsUseCase - GetContacts - GetContactsByUserUUID: %w", err)
	}
	return contacts, nil
}

func (uc *ContactsUseCase) AddContact(ctx context.Context, contactUserName string, userUuid string) error {
	// Check username exists in user data repository
	contactUserUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, contactUserName)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - AddContacts - GetUserUUIDByUsername: %w", err)
	}

	if contactUserUUID == nil {
		return entity.ErrUserNameNotFound
	}

	// Check contact exist from contacts data repository
	exist, err := uc.repo.CheckContactExist(ctx, userUuid, *contactUserUUID)
	if err != nil {
		return err
	}

	if exist {
		// If contact exist, then update remove status from data respository
		err := uc.repo.UpdateRemovedStatus(ctx, entity.ContactsDTO{
			UserUUID:        userUuid,
			ContactUserUUID: *contactUserUUID,
			Removed:         false,
		})

		if err != nil {
			return fmt.Errorf("ContactsUseCase - AddContacts - uc.repo.UpdateRemoved: %w", err)
		}
		return nil
	}

	// Convert into contactsDTO
	contactsDTO := entity.ContactsDTO{
		UserUUID:         userUuid,
		ContactUserUUID:  *contactUserUUID,
		ConversationUUID: uuid.New().String(),
	}

	// Store contact into contacts data repository
	err = uc.repo.StoreContacts(ctx, contactsDTO)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - AddContacts - uc.repo.StoreContacts: %w", err)
	}
	return nil
}

func (uc *ContactsUseCase) RemoveContact(ctx context.Context, contactUserName string, userUuid string) error {
	// Check username exists from user data repository
	contactUserUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, contactUserName)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - RemoveContact - GetUserUUIDByUsername: %w", err)
	}

	if contactUserUUID == nil {
		return entity.ErrUserNameNotFound
	}

	// Check contact exist from data repository
	exist, err := uc.repo.CheckContactExist(ctx, userUuid, *contactUserUUID)
	if err != nil {
		return err
	}

	// Return error if contact does not exist
	if !exist {
		return entity.ErrContactDoesNotExists
	}

	// Convert into contactsDTO
	contactsDTO := entity.ContactsDTO{
		UserUUID:        userUuid,
		ContactUserUUID: *contactUserUUID,
		Removed:         true,
	}

	// Set remove status in contacts data repository
	err = uc.repo.UpdateRemovedStatus(ctx, contactsDTO)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - RemoveContact - uc.repo.UpdateRemoved: %w", err)
	}
	return nil
}

func (uc *ContactsUseCase) UpdateBlockContact(ctx context.Context, contactUserName string, userUuid string, block bool) error {
	// Check username exists in user data repository
	contactUserUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, contactUserName)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - UpdateBlockContact - GetUserUUIDByUsername: %w", err)
	}

	// Return error if username does not exist
	if contactUserUUID == nil {
		return entity.ErrUserNameNotFound
	}

	// Check contact exist from contacts data repository
	exist, err := uc.repo.CheckContactExist(ctx, userUuid, *contactUserUUID)
	if err != nil {
		return err
	}

	// Return error if contact does not exist
	if !exist {
		return entity.ErrContactDoesNotExists
	}

	// Convert into contactsDTO
	contactsDTO := entity.ContactsDTO{
		UserUUID:        userUuid,
		ContactUserUUID: *contactUserUUID,
		Blocked:         block,
	}
	// Update 'blocked' column in contacts data respository
	err = uc.repo.UpdateBlockedStatus(ctx, contactsDTO)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - UpdateBlockContact - uc.repo.UpdateBlockedStatus: %w", err)
	}
	return nil
}
