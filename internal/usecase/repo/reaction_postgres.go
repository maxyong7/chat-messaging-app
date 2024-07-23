package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// ReactionRepo -.
type ReactionRepo struct {
	// *postgres.Postgres
	*sql.DB
}

// New -.
func NewReaction(pg *sql.DB) *ReactionRepo {
	return &ReactionRepo{pg}
}

// StoreReaction -.
func (r *ReactionRepo) StoreReaction(ctx context.Context, srDTO entity.StoreReactionDTO) error {
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ReactionRepo - StoreReaction - failed to begin transaction: %w", err)
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

	upsertReactionSQL := `
		INSERT INTO reaction (message_uuid, user_uuid, reaction_type)
		VALUES ($1, $2, $3)
		ON CONFLICT (message_uuid, user_uuid)
		DO UPDATE SET
			reaction_type = EXCLUDED.reaction_type
		`
	_, err = tx.ExecContext(ctx, upsertReactionSQL, srDTO.MessageUUID, srDTO.SenderUUID, srDTO.ReactionType)
	if err != nil {
		return fmt.Errorf("failed to execute insert upsertReactionSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ReactionRepo - StoreReaction - failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ReactionRepo) GetReactions(ctx context.Context, messageUUID string) ([]entity.GetReaction, error) {
	getReactionSQL := `
		SELECT
			r.reaction_type,
			ui.first_name,
			ui.last_name,
			ui.avatar
		FROM reaction r
		LEFT JOIN user_info ui ON r.user_uuid = ui.user_uuid
		WHERE r.message_uuid = $1
	`

	// Execute the final query.
	rows, err := r.QueryContext(ctx, getReactionSQL, messageUUID)
	if err != nil {
		fmt.Println("GetReactions - getReactionSQL err: ", err)
		return nil, err
	}
	defer rows.Close()

	// Process the results.
	var reactions []entity.GetReaction
	for rows.Next() {
		var reaction entity.GetReaction
		if err := rows.Scan(
			&reaction.ReactionType,
			&reaction.FirstName,
			&reaction.LastName,
			&reaction.Avatar,
		); err != nil {
			fmt.Println("GetReaction - rows.Scan err: ", err)
			return nil, err
		}
		reactions = append(reactions, reaction)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("GetReactions - rows.Err(): ", err)
	}

	return reactions, nil
}

func (r *ReactionRepo) RemoveReaction(ctx context.Context, rr entity.RemoveReactionDTO) error {
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ReactionRepo - RemoveReaction - failed to begin transaction: %w", err)
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

	deleteReactionSQL := `
		DELETE FROM reaction
		WHERE message_uuid = $1
		AND user_uuid = $2
		`
	_, err = tx.ExecContext(ctx, deleteReactionSQL, rr.MessageUUID, rr.SenderUUID)
	if err != nil {
		return fmt.Errorf("failed to execute insert deleteReactionSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ReactionRepo - RemoveReaction - failed to commit transaction: %w", err)
	}

	return nil
}
