package repo

import (
	"context"
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/pkg/postgres"
)

// ConversationRepo -.
type ConversationRepo struct {
	*postgres.Postgres
}

// New -.
func NewConversation(pg *postgres.Postgres) *ConversationRepo {
	return &ConversationRepo{pg}
}

// GetUserInfo -.
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

	rows, err := r.Pool.Query(ctx, query, reqParam.UserID)
	if err != nil {
		return nil, fmt.Errorf("ConversationRepo - GetConversation - r.Pool.Query: %w", err)
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
		// return nil, entity.ErrNoConversationFound
		return nil, nil
	}

	finalQuery := `
		SELECT 
			c.last_message, 
			c.last_sent_user_uuid, 
			c.title, 
			c.last_message_created_at, 
			c.conversation_type,
			ui.first_name,
			ui.last_name,
			ui.avatar
		FROM conversations c
		WHERE conversation_uuid = IN($1)
		AND last_message_created_at < $2
		LEFT JOIN user_info ui ON conversations.last_sent_user_uuid = user_info.user_uuid
		ORDER BY last_message_created_at DESC
		LIMIT $3;
	`

	// Execute the final query.
	rows, err = r.Pool.Query(ctx, finalQuery, pg.In(conversationUUIDs), reqParam.Cursor, reqParam.Limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Process the results.
	var conversations []entity.Conversations
	for rows.Next() {
		var conv entity.Conversations
		if err := rows.Scan(
			&conv.LastMessage,
			&conv.LastSentUser,
			&conv.Title,
			&conv.LastMessageCreatedAt,
			&conv.Type,
			&conv.LastSentUser.FirstName,
			&conv.LastSentUser.LastName,
			&conv.LastSentUser.Avatar,
		); err != nil {
			log.Fatal(err)
		}
		conversations = append(conversations, conv)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return conversations, nil
}
