package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
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

func (r *MessageRepo) GetMessages(ctx context.Context, reqParam entity.RequestParams, conversationUUID string) ([]entity.GetMessageDTO, error) {
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
		LEFT JOIN reaction r ON m.message_uuid = r.message_uuid
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

func (r *MessageRepo) ValidateMessageSentByUser(msgReq usecase.MessageRequest) (bool, error) {
	validateMessageSQL := `
		SELECT 1
		FROM messages 
		WHERE (message_uuid = $1 AND user_uuid = $2)
	`

	var exists int
	err := r.QueryRow(validateMessageSQL, msgReq.Data.MessageUUID, msgReq.Data.SenderUUID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("UserInfoRepo - CheckUserExist -  r.QueryRowContext: %w", err)
	}

	if exists > 0 {
		return true, nil
	}

	return false, nil

}

func (r *MessageRepo) DeleteMessage(msgReq usecase.MessageRequest) error {
	ctx := context.Background()
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
	_, err = tx.ExecContext(ctx, deleteMessageSQL, msgReq.Data.MessageUUID, msgReq.Data.SenderUUID)
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
