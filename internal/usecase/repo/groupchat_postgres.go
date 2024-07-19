package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// GroupChatRepo -.
type GroupChatRepo struct {
	*sql.DB
}

// New -.
func NewGroupChat(pg *sql.DB) *GroupChatRepo {
	return &GroupChatRepo{pg}
}

// CreateGroupChat -.
func (r *GroupChatRepo) CreateGroupChat(ctx context.Context, groupChat entity.GroupChatRequest) error {
	conversationUUID := uuid.New().String()
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("GroupChatRepo - CreateGroupChat - failed to begin transaction: %w", err)
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

	// Insert row for user
	insertGroupChatSQL := `
	INSERT INTO participants (user_uuid, conversation_uuid)
	VALUES ($1, $2)
	`
	_, err = tx.ExecContext(ctx, insertGroupChatSQL, groupChat.UserUUID, conversationUUID)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertGroupChatSQL query for user: %w", err)
	}

	// Insert rows for participants
	for _, participant := range groupChat.Participants {
		_, err = tx.ExecContext(ctx, insertGroupChatSQL, participant, conversationUUID)
		if err != nil {
			return fmt.Errorf("failed to execute insert insertGroupChatSQL query for participants: %w", err)
		}
	}

	// Insert row in conversations table
	insertConversationsSQL := `
	INSERT INTO conversations (
		conversation_uuid, title, conversation_type
	) VALUES ($1, $2, $3)
	`
	_, err = tx.ExecContext(ctx, insertConversationsSQL, conversationUUID, groupChat.Title, entity.GroupMessageConversationType)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertConversationsSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("GroupChatRepo - CreateGroupChat - failed to commit transaction: %w", err)
	}

	return nil
}

// AddParticipants -.
func (r *GroupChatRepo) AddParticipants(ctx context.Context, groupChat entity.GroupChatRequest) error {
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("GroupChatRepo - AddParticipants - failed to begin transaction: %w", err)
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

	// Insert rows for participants
	for _, participant := range groupChat.Participants {
		addParticipantsSQL := `
		INSERT INTO participants (user_uuid, conversation_uuid)
		VALUES ($1, $2)
		`
		_, err = tx.ExecContext(ctx, addParticipantsSQL, participant, groupChat.ConversationUUID)
		if err != nil {
			return fmt.Errorf("failed to execute insert addParticipantsSQL query for participants: %w", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("GroupChatRepo - AddParticipants - failed to commit transaction: %w", err)
	}

	return nil
}

// RemoveParticipant -.
func (r *GroupChatRepo) RemoveParticipants(ctx context.Context, groupChat entity.GroupChatRequest) error {
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("GroupChatRepo - RemoveParticipant - failed to begin transaction: %w", err)
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

	// Insert rows for participants
	for _, participant := range groupChat.Participants {
		removeParticipantsSQL := `
		DELETE FROM participants 
		WHERE user_uuid = $1 
		AND conversation_uuid = $2
		`
		_, err = tx.ExecContext(ctx, removeParticipantsSQL, participant, groupChat.ConversationUUID)
		if err != nil {
			return fmt.Errorf("failed to execute insert removeParticipantsSQL query for participants: %w", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("GroupChatRepo - RemoveParticipants - failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateGroupTitle -.
func (r *GroupChatRepo) UpdateGroupTitle(ctx context.Context, groupChat entity.GroupChatRequest) error {
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("GroupChatRepo - UpdateGroupTitle - failed to begin transaction: %w", err)
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

	UpdateGroupTitleSQL := `
		UPDATE conversations 
		SET title = $1
		WHERE conversation_uuid = $2
		`
	_, err = tx.ExecContext(ctx, UpdateGroupTitleSQL, groupChat.Title, groupChat.ConversationUUID)
	if err != nil {
		return fmt.Errorf("failed to execute insert UpdateGroupTitleSQL query for participants: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("GroupChatRepo - UpdateGroupTitle - failed to commit transaction: %w", err)
	}

	return nil
}

// ValidateUserInGroupChat -.
func (r *GroupChatRepo) ValidateUserInGroupChat(ctx context.Context, conversationUUID string, userUUID string) (bool, error) {
	// Check if the user already exists
	validateUserInGroupChatSQL := `
	SELECT 1 
	FROM participants
	WHERE (conversation_uuid = $1 or user_uuid = $2)
	`

	var exists int
	err := r.QueryRowContext(ctx, validateUserInGroupChatSQL, conversationUUID, userUUID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("GroupChatRepo - ValidateUserInGroupChat -  r.QueryRowContext: %w", err)
	}

	if exists > 0 {
		return true, nil
	}

	return false, nil
}
