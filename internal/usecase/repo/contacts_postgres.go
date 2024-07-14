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

// // GetUserInfo -.
// func (r *ContactsRepo) GetConversations(ctx context.Context, reqParam entity.RequestParams) ([]entity.Conversations, error) {
// 	// Define the SQL query.
// 	query := `
// 		WITH combined_conversations AS (
// 			SELECT conversation_uuid
// 			FROM contacts
// 			WHERE user_uuid = $1
// 			UNION
// 			SELECT conversation_uuid
// 			FROM participants
// 			WHERE user_uuid = $1
// 		)
// 		SELECT conversation_uuid
// 		FROM combined_conversations
// 		`

// 	rows, err := r.QueryContext(ctx, query, reqParam.UserID)
// 	if err != nil {
// 		return nil, fmt.Errorf("ContactsRepo - GetConversation - r.QueryContext: %w", err)
// 	}
// 	defer rows.Close()

// 	// Collect the conversation UUIDs.
// 	var conversationUUIDs []string
// 	for rows.Next() {
// 		var conversationUUID string
// 		if err := rows.Scan(&conversationUUID); err != nil {
// 			return nil, fmt.Errorf("ContactsRepo - GetConversation - rows.Scan: %w", err)
// 		}
// 		conversationUUIDs = append(conversationUUIDs, conversationUUID)
// 	}

// 	// If there are no conversation UUIDs, return early.
// 	if len(conversationUUIDs) == 0 {
// 		fmt.Println("No conversations found.")
// 		// return nil, entity.ErrNoConversationFound
// 		return nil, nil
// 	}

// 	finalQuery := `
// 		SELECT
// 			c.last_message,
// 			c.last_sent_user_uuid,
// 			c.title,
// 			c.last_message_created_at,
// 			c.conversation_type,
// 			ui.first_name,
// 			ui.last_name,
// 			ui.avatar
// 		FROM conversations c
// 		WHERE conversation_uuid = IN($1)
// 		AND last_message_created_at < $2
// 		LEFT JOIN user_info ui ON conversations.last_sent_user_uuid = user_info.user_uuid
// 		ORDER BY last_message_created_at DESC
// 		LIMIT $3;
// 	`

// 	// Execute the final query.
// 	rows, err = r.QueryContext(ctx, finalQuery, pg.In(conversationUUIDs), reqParam.Cursor, reqParam.Limit)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	// Process the results.
// 	var conversations []entity.Conversations
// 	for rows.Next() {
// 		var conv entity.Conversations
// 		if err := rows.Scan(
// 			&conv.LastMessage,
// 			&conv.LastSentUser,
// 			&conv.Title,
// 			&conv.LastMessageCreatedAt,
// 			&conv.Type,
// 			&conv.LastSentUser.FirstName,
// 			&conv.LastSentUser.LastName,
// 			&conv.LastSentUser.Avatar,
// 		); err != nil {
// 			log.Fatal(err)
// 		}
// 		conversations = append(conversations, conv)
// 	}
// 	if err := rows.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	return conversations, nil
// }

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
			return nil, fmt.Errorf("ConversationRepo - GetConversation - rows.Scan: %w", err)
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

// StoreConversation -.
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

	insertMessagesSQL := `
		INSERT INTO contacts (user_uuid, contact_user_uuid, conversation_uuid, blocked)
		VALUES ($1, $2, $3, $4)
		`
	_, err = tx.ExecContext(ctx, insertMessagesSQL, contacts.UserUUID, contacts.ContactUserUUID, contacts.ConversationUUID, contacts.Blocked)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertMessagesSQL query: %w", err)
	}

	insertMessagesAsContactSQL := `
		INSERT INTO contacts (user_uuid, contact_user_uuid, conversation_uuid, blocked)
		VALUES ($1, $2, $3, $4)
		`
	_, err = tx.ExecContext(ctx, insertMessagesAsContactSQL, contacts.ContactUserUUID, contacts.UserUUID, contacts.ConversationUUID, contacts.Blocked)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertMessagesAsContactSQL query: %w", err)
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
		return fmt.Errorf("ContactsRepo - StoreConversation - failed to commit transaction: %w", err)
	}

	return nil
}
