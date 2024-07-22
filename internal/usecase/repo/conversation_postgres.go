package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
)

// ConversationRepo -.
type ConversationRepo struct {
	// *postgres.Postgres
	*sql.DB
}

// New -.
func NewConversation(pg *sql.DB) *ConversationRepo {
	return &ConversationRepo{pg}
}

func (r *ConversationRepo) GetConversations(ctx context.Context, reqParam entity.RequestParams) ([]entity.Conversations, error) {
	// Define the SQL query.
	query := `
		WITH combined_conversations AS (
			SELECT conversation_uuid
			FROM contacts
			WHERE user_uuid = $1
			UNION
			SELECT conversation_uuid
			FROM participants
			WHERE user_uuid = $1
		)
		SELECT conversation_uuid
		FROM combined_conversations
		`

	rows, err := r.QueryContext(ctx, query, reqParam.UserID)
	if err != nil {
		return nil, fmt.Errorf("ConversationRepo - GetConversation - r.QueryContext: %w", err)
	}
	defer rows.Close()

	// Collect the conversation UUIDs.
	var conversationUUIDs []string
	for rows.Next() {
		var conversationUUID string
		if err := rows.Scan(&conversationUUID); err != nil {
			return nil, fmt.Errorf("ConversationRepo - GetConversation - rows.Scan: %w", err)
		}
		conversationUUIDs = append(conversationUUIDs, conversationUUID)
	}

	// If there are no conversation UUIDs, return early.
	if len(conversationUUIDs) == 0 {
		fmt.Println("No conversations found.")
		return nil, nil
	}

	finalQuery := `
		SELECT
			c.last_message,
			c.title,
			c.last_message_created_at,
			c.conversation_type,
			ui.first_name,
			ui.last_name,
			ui.avatar
		FROM conversations c
		LEFT JOIN user_info ui ON c.last_sent_user_uuid = ui.user_uuid
		WHERE c.conversation_uuid = ANY($1)
		AND c.last_message IS NOT NULL			
		AND c.last_message_created_at < $2
		ORDER BY c.last_message_created_at DESC
		LIMIT $3;
	`

	// Execute the final query.
	rows, err = r.QueryContext(ctx, finalQuery, pq.Array(conversationUUIDs), reqParam.Cursor, reqParam.Limit)
	if err != nil {
		fmt.Println("GetConversations - finalQuery err: ", err)
		return nil, err
	}
	defer rows.Close()

	// Process the results.
	var conversations []entity.Conversations
	for rows.Next() {
		var conv entity.Conversations
		if err := rows.Scan(
			&conv.LastMessage,
			&conv.Title,
			&conv.LastMessageCreatedAt,
			&conv.Type,
			&conv.LastSentUser.FirstName,
			&conv.LastSentUser.LastName,
			&conv.LastSentUser.Avatar,
		); err != nil {
			fmt.Println("GetConversations - rows.Scan err: ", err)
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("GetConversations - rows.Err(): ", err)
	}

	return conversations, nil
}

// StoreConversation -.
func (r *ConversationRepo) StoreConversation(msg usecase.MessageRequest) error {
	ctx := context.Background()
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ConversationRepo - StoreConversation - failed to begin transaction: %w", err)
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
		INSERT INTO messages (message_uuid, conversation_uuid, user_uuid, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
		`
	_, err = tx.ExecContext(ctx, insertMessagesSQL, msg.Data.MessageUUID, msg.Data.ConversationUUID, msg.Data.SenderUUID, msg.Data.Content, msg.Data.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertMessagesSQL query: %w", err)
	}

	upsertConversationsSQL := `
	INSERT INTO conversations (
		conversation_uuid,
		last_message,
		last_sent_user_uuid,
		last_message_created_at
	) VALUES ($1, $2, $3, $4)
	ON CONFLICT (conversation_uuid) 
	DO UPDATE SET
		last_message = EXCLUDED.last_message,
		last_sent_user_uuid = EXCLUDED.last_sent_user_uuid,
		title = EXCLUDED.title,
		last_message_created_at = EXCLUDED.last_message_created_at
	`
	_, err = tx.ExecContext(ctx, upsertConversationsSQL, msg.Data.ConversationUUID, msg.Data.Content, msg.Data.SenderUUID, msg.Data.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute insert upsertConversationsSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ConversationRepo - StoreConversation - failed to commit transaction: %w", err)
	}

	return nil
}
