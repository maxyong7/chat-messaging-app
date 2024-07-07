package repo

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/pkg/postgres"
)

// VerificationRepo -.
type VerificationRepo struct {
	*postgres.Postgres
}

// New -.
func NewVerification(pg *postgres.Postgres) *VerificationRepo {
	return &VerificationRepo{pg}
}

// GetUserInfo -.
func (r *VerificationRepo) GetUserInfo(ctx context.Context, verification entity.Verification) (*entity.UserInfoDTO, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("users").
		Where(
			squirrel.Or{
				squirrel.Eq{"username": verification.Username},
				squirrel.Eq{"email": verification.Email},
			},
		).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("VerificationRepo - GetUserInfo - r.Builder: %w", err)
	}

	var userInfo entity.UserInfoDTO
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Password, &userInfo.Email)
	if err != nil {
		if err != pgx.ErrNoRows {
			return &userInfo, nil
		}
		return nil, fmt.Errorf("VerificationRepo - GetUserInfo - r.Pool.QueryRow: %w", err)
	}

	return &userInfo, nil
}

// // Store -.
// func (r *VerificationRepo) Store(ctx context.Context, t entity.UserInfoDTO) error {
// 	sql, args, err := r.Builder.
// 		Insert("history").
// 		Columns("source, destination, original, translation").
// 		Values(t.Source, t.Destination, t.Original, t.Translation).
// 		ToSql()
// 	if err != nil {
// 		return fmt.Errorf("VerificationRepo - Store - r.Builder: %w", err)
// 	}

// 	_, err = r.Pool.Exec(ctx, sql, args...)
// 	if err != nil {
// 		return fmt.Errorf("VerificationRepo - Store - r.Pool.Exec: %w", err)
// 	}

// 	return nil
// }
