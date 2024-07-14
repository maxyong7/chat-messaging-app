package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// ContactsRepo -.
type ContactsRepo struct {
	// *postgres.Postgres
	*sql.DB
}

// New -.
func NewContacts(pg *sql.DB) *ContactsRepo {
	return &ContactsRepo{pg}
}

func (r *ContactsRepo) CheckContactExist(ctx context.Context, userUuid string, contactUserUuid string) (bool, error) {
	// Define the SQL query.
	query := `
		SELECT 1 FROM contacts
		WHERE (user_uuid = $1 AND contact_user_uuid = $2)
		`

	var exists int
	err := r.QueryRowContext(ctx, query, userUuid, contactUserUuid).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("UserInfoRepo - CheckUserExist -  r.QueryRowContext: %w", err)
	}

	if exists > 0 {
		return true, nil
	}

	return false, nil

}

func (r *ContactsRepo) GetContactsByUserUUID(ctx context.Context, userUuid string) ([]entity.Contacts, error) {
	// Define the SQL query.
	query := `
		SELECT 
			ui.first_name,
			ui.last_name,
			ui.avatar,
			ct.conversation_uuid,
			ct.blocked
		FROM contacts ct
		LEFT JOIN user_info ui ON ct.contact_user_uuid = ui.user_uuid
		WHERE ct.user_uuid = $1
		AND ct.removed != true
		ORDER BY ui.first_name;
		`

	rows, err := r.QueryContext(ctx, query, userUuid)
	if err != nil {
		return nil, fmt.Errorf("ContactsRepo - GetContactsByUserUUID - r.QueryContext: %w", err)
	}
	defer rows.Close()
	// Collect the conversation UUIDs.
	var contacts []entity.Contacts
	for rows.Next() {
		var contact entity.Contacts
		if err := rows.Scan(&contact.FirstName, &contact.LastName, &contact.Avatar, &contact.ConversationUUID, &contact.Blocked); err != nil {
			return nil, fmt.Errorf("ContactsRepo - GetContactsByUserUUID - rows.Scan: %w", err)
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

// StoreContacts -.
func (r *ContactsRepo) StoreContacts(ctx context.Context, contacts entity.ContactsDTO) error {
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ContactsRepo - StoreContacts - failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back if it doesn't commit
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; rollback
		}
	}()

	insertContactAsUserSQL := `
		INSERT INTO contacts (user_uuid, contact_user_uuid, conversation_uuid, blocked)
		VALUES ($1, $2, $3, $4)
		`
	_, err = tx.ExecContext(ctx, insertContactAsUserSQL, contacts.UserUUID, contacts.ContactUserUUID, contacts.ConversationUUID, contacts.Blocked)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertContactAsUserSQL query: %w", err)
	}

	insertContactAsContactSQL := `
		INSERT INTO contacts (user_uuid, contact_user_uuid, conversation_uuid, blocked)
		VALUES ($1, $2, $3, $4)
		`
	_, err = tx.ExecContext(ctx, insertContactAsContactSQL, contacts.ContactUserUUID, contacts.UserUUID, contacts.ConversationUUID, contacts.Blocked)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertContactAsContactSQL query: %w", err)
	}

	insertConversationsSQL := `
	INSERT INTO conversations (
		conversation_uuid
	) VALUES ($1)
	`
	_, err = tx.ExecContext(ctx, insertConversationsSQL, contacts.ConversationUUID)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertConversationsSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ContactsRepo - StoreContacts - failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateBlocked -.
func (r *ContactsRepo) UpdateBlocked(ctx context.Context, contacts entity.ContactsDTO) error {
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ContactsRepo - UpdateBlocked - failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back if it doesn't commit
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; rollback
		}
	}()

	updateBlockedSQL := `
		UPDATE contacts 
		SET blocked = $1
		WHERE user_uuid = $2
		AND contact_user_uuid = $3
		`
	_, err = tx.ExecContext(ctx, updateBlockedSQL, contacts.Blocked, contacts.UserUUID, contacts.ContactUserUUID)
	if err != nil {
		return fmt.Errorf("failed to execute insert updateBlockedSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ContactsRepo - UpdateBlocked - failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateRemoved -.
func (r *ContactsRepo) UpdateRemoved(ctx context.Context, contacts entity.ContactsDTO) error {
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ContactsRepo - UpdateRemoved - failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back if it doesn't commit
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; rollback
		}
	}()

	updateRemovedSQL := `
		UPDATE contacts 
		SET removed = $1
		WHERE user_uuid = $2
		AND contact_user_uuid = $3
		`
	_, err = tx.ExecContext(ctx, updateRemovedSQL, contacts.Removed, contacts.UserUUID, contacts.ContactUserUUID)
	if err != nil {
		return fmt.Errorf("failed to execute insert updateRemovedSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ContactsRepo - UpdateRemoved - failed to commit transaction: %w", err)
	}

	return nil
}
