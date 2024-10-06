package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// MessageRepo -.
type MessageRepo struct {
	// *postgres.Postgres
	*sql.DB
}

// New -.
func NewMessage(pg *sql.DB) *MessageRepo {
	return &MessageRepo{pg}
}

func (r *MessageRepo) GetMessages(ctx context.Context, reqParam entity.RequestParamsDTO, conversationUUID string) ([]entity.GetMessageDTO, error) {
	getMessagesSQL := `
		SELECT
			m.message_uuid,
			m.user_uuid,
			m.content,
			m.created_at,
			ui.first_name,
			ui.last_name,
			ui.avatar
		FROM messages m
		LEFT JOIN user_info ui ON m.user_uuid = ui.user_uuid
		WHERE m.conversation_uuid = $1
		AND m.created_at < $2
		ORDER BY m.created_at DESC
		LIMIT $3;
	`

	// Execute the final query.
	rows, err := r.QueryContext(ctx, getMessagesSQL, conversationUUID, reqParam.Cursor, reqParam.Limit)
	if err != nil {
		fmt.Println("GetMessages - getMessagesSQL err: ", err)
		return nil, err
	}
	defer rows.Close()

	// Process the results.
	var messages []entity.GetMessageDTO
	for rows.Next() {
		var msg entity.GetMessageDTO
		if err := rows.Scan(
			&msg.MessageUUID,
			&msg.User.UserUUID,
			&msg.Content,
			&msg.CreatedAt,
			&msg.User.FirstName,
			&msg.User.LastName,
			&msg.User.Avatar,
		); err != nil {
			fmt.Println("GetConversations - rows.Scan err: ", err)
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("GetConversations - rows.Err(): ", err)
	}

	return messages, nil
}

func (r *MessageRepo) ValidateMessageSentByUser(ctx context.Context, msg entity.MessageDTO) (bool, error) {
	validateMessageSQL := `
		SELECT 1
		FROM messages 
		WHERE (message_uuid = $1 AND user_uuid = $2)
	`

	var exists int
	err := r.QueryRowContext(ctx, validateMessageSQL, msg.MessageUUID, msg.UserUUID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("MessageRepo - ValidateMessageSentByUser -  validateMessageSQL: %w", err)
	}

	if exists > 0 {
		return true, nil
	}

	return false, nil

}

func (r *MessageRepo) DeleteMessage(ctx context.Context, msg entity.MessageDTO) error {
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("MessageRepo - DeleteMessage - failed to begin transaction: %w", err)
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

	deleteMessageSQL := `
		DELETE FROM messages
		WHERE message_uuid = $1
		AND user_uuid = $2
		`
	_, err = tx.ExecContext(ctx, deleteMessageSQL, msg.MessageUUID, msg.UserUUID)
	if err != nil {
		return fmt.Errorf("failed to execute insert deleteMessageSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("MessageRepo - DeleteMessage - failed to commit transaction: %w", err)
	}

	return nil
}

func (r *MessageRepo) UpdateSeenStatus(ctx context.Context, seenStatus entity.SeenStatusDTO) error {
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("MessageRepo - UpdateSeenStatus - failed to begin transaction: %w", err)
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

	insertSeenStatusSQL := `
		INSERT INTO seen_status (message_uuid, user_uuid, seen_timestamp)
		SELECT m.message_uuid, $1, NOW()
		FROM messages m
		LEFT JOIN seen_status s
		ON m.message_uuid = s.message_uuid
		WHERE m.conversation_uuid = $2
		AND m.user_uuid <> $1
		AND s.message_uuid IS NULL;
		`
	_, err = tx.ExecContext(ctx, insertSeenStatusSQL, seenStatus.UserUUID, seenStatus.ConversationUUID)
	if err != nil {
		return fmt.Errorf("failed to execute insert deleteMessageSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("MessageRepo - UpdateSeenStatus - failed to commit transaction: %w", err)
	}

	return nil
}

func (r *MessageRepo) GetSeenStatus(ctx context.Context, messageUUID string) ([]entity.GetSeenStatusDTO, error) {
	getSeenStatusSQL := `
		SELECT
			s.seen_timestamp,
			ui.first_name,
			ui.last_name,
			ui.avatar
		FROM seen_status s
		LEFT JOIN user_info ui ON s.user_uuid = ui.user_uuid
		WHERE s.message_uuid = $1
		ORDER BY s.seen_timestamp DESC
	`

	// Execute the final query.
	rows, err := r.QueryContext(ctx, getSeenStatusSQL, messageUUID)
	if err != nil {
		fmt.Println("GetSeenStatus - getSeenStatusSQL err: ", err)
		return nil, err
	}
	defer rows.Close()

	// Process the results.
	var seenStatuses []entity.GetSeenStatusDTO
	for rows.Next() {
		var seenStatus entity.GetSeenStatusDTO
		if err := rows.Scan(
			&seenStatus.SeenTimestamp,
			&seenStatus.FirstName,
			&seenStatus.LastName,
			&seenStatus.Avatar,
		); err != nil {
			fmt.Println("GetSeenStatus - rows.Scan err: ", err)
			return nil, err
		}
		seenStatuses = append(seenStatuses, seenStatus)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("GetSeenStatus - rows.Err(): ", err)
	}

	return seenStatuses, nil
}

func (r *MessageRepo) SearchMessage(ctx context.Context, keyword string, conversationUUID string) ([]entity.SearchMessageDTO, error) {
	getMessagesSQL := `
		SELECT
			m.message_uuid,
			m.user_uuid,
			m.content,
			m.created_at,
			ui.first_name,
			ui.last_name,
			ui.avatar
		FROM messages m
		LEFT JOIN user_info ui ON m.user_uuid = ui.user_uuid
		WHERE m.content ILIKE '%' || $1 || '%' 
		AND m.conversation_uuid = $2
	`

	// Execute the final query.
	rows, err := r.QueryContext(ctx, getMessagesSQL, keyword, conversationUUID)
	if err != nil {
		fmt.Println("GetMessages - getMessagesSQL err: ", err)
		return nil, err
	}
	defer rows.Close()

	// Process the results.
	var messages []entity.SearchMessageDTO
	for rows.Next() {
		var msg entity.SearchMessageDTO
		if err := rows.Scan(
			&msg.MessageUUID,
			&msg.User.UserUUID,
			&msg.Content,
			&msg.CreatedAt,
			&msg.User.FirstName,
			&msg.User.LastName,
			&msg.User.Avatar,
		); err != nil {
			fmt.Println("GetConversations - rows.Scan err: ", err)
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("GetConversations - rows.Err(): ", err)
	}

	return messages, nil
}
