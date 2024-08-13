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
	// Query 'contacts' table from contacts data repository
	// Then join 'user_info' table on 'user_uuid' to get user's firstname, lastname and avatar
	contacts, err := uc.repo.GetContactsByUserUUID(ctx, userUuid)
	if err != nil {
		return nil, fmt.Errorf("ContactsUseCase - GetContacts - GetContactsByUserUUID: %w", err)
	}
	return contacts, nil
}

func (uc *ContactsUseCase) AddContact(ctx context.Context, contactUserName string, userUuid string) error {
	// Check username exists in user repository by querying 'user_credentials' table
	contactUserUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, contactUserName)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - AddContacts - GetUserUUIDByUsername: %w", err)
	}

	// Return error if contact is not found. Will be handled by controller
	if contactUserUUID == nil {
		return entity.ErrUserNameNotFound
	}

	// Check contact exist from contacts data repository by querying 'contacts' table
	exist, err := uc.repo.CheckContactExist(ctx, userUuid, *contactUserUUID)
	if err != nil {
		return err
	}

	if exist {
		// If contact exist, then update 'remove' column on 'contacts' table in contacts data respository
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

	// Convert arguments into contactsDTO
	contactsDTO := entity.ContactsDTO{
		UserUUID:         userUuid,
		ContactUserUUID:  *contactUserUUID,
		ConversationUUID: uuid.New().String(),
	}

	// Store contact into 'contacts' table in contacts data repository
	err = uc.repo.StoreContacts(ctx, contactsDTO)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - AddContacts - uc.repo.StoreContacts: %w", err)
	}
	return nil
}

func (uc *ContactsUseCase) RemoveContact(ctx context.Context, contactUserName string, userUuid string) error {
	// Check if user exists by querying 'user_credentials' table in user data repository
	contactUserUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, contactUserName)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - RemoveContact - GetUserUUIDByUsername: %w", err)
	}

	// Return error if username not found. Will be handled by controller
	if contactUserUUID == nil {
		return entity.ErrUserNameNotFound
	}

	// Check contact exist from contacts data repository by querying 'contacts' table
	exist, err := uc.repo.CheckContactExist(ctx, userUuid, *contactUserUUID)
	if err != nil {
		return err
	}

	// Return error if contact does not exist. Will be handled by controller
	if !exist {
		return entity.ErrContactDoesNotExists
	}

	// Convert arguments into contactsDTO
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
	// Check if user exists by querying 'user_credentials' table in user data repository
	contactUserUUID, err := uc.userInfoRepo.GetUserUUIDByUsername(ctx, contactUserName)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - UpdateBlockContact - GetUserUUIDByUsername: %w", err)
	}

	// Return error if username does not exist. Will be handled by controller
	if contactUserUUID == nil {
		return entity.ErrUserNameNotFound
	}

	// Check contact exist from contacts data repository by querying 'contacts' table
	exist, err := uc.repo.CheckContactExist(ctx, userUuid, *contactUserUUID)
	if err != nil {
		return err
	}

	// Return error if contact does not exist. Will be handled by controller
	if !exist {
		return entity.ErrContactDoesNotExists
	}

	// Convert arguments into contactsDTO
	contactsDTO := entity.ContactsDTO{
		UserUUID:        userUuid,
		ContactUserUUID: *contactUserUUID,
		Blocked:         block,
	}
	// Update 'blocked' column in 'contacts' table from contacts data repository
	err = uc.repo.UpdateBlockedStatus(ctx, contactsDTO)
	if err != nil {
		return fmt.Errorf("ContactsUseCase - UpdateBlockContact - uc.repo.UpdateBlockedStatus: %w", err)
	}
	return nil
}
