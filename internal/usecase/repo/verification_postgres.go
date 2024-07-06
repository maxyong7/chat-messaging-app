package repo

import (
	"context"
	"fmt"

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
	sql, _, err := r.Builder.
		Select("*").
		From("user").
		Where("username", verification.Username).
		Where("email", verification.Email).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("VerificationRepo - GetUserInfo - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("VerificationRepo - GetUserInfo - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := entity.UserInfoDTO{}
	err = rows.Scan(&entities.Email, &entities.Username, &entities.Password)
	if err != nil {
		return nil, fmt.Errorf("VerificationRepo - GetUserInfo - rows.Scan: %w", err)
	}

	return &entities, nil
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
