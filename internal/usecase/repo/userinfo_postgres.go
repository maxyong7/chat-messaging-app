package repo

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
	"github.com/maxyong7/chat-messaging-app/pkg/postgres"
)

// UserInfoRepo -.
type UserInfoRepo struct {
	*postgres.Postgres
}

// New -.
func NewUserInfo(pg *postgres.Postgres) *UserInfoRepo {
	return &UserInfoRepo{pg}
}

// GetUserInfo -.
func (r *UserInfoRepo) GetUserInfo(ctx context.Context, userInfo entity.UserInfo) (*entity.UserInfoDTO, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("user_credentials").
		Where(
			squirrel.Or{
				squirrel.Eq{"username": userInfo.Username},
				squirrel.Eq{"email": userInfo.Email},
			},
		).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("UserInfoRepo - GetUserInfo - r.Builder: %w", err)
	}

	var userInfoDTO entity.UserInfoDTO
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&userInfoDTO.ID, &userInfoDTO.Username, &userInfoDTO.Password, &userInfoDTO.Email)
	if err != nil {
		if err != pgx.ErrNoRows {
			return &userInfoDTO, nil
		}
		return nil, fmt.Errorf("UserInfoRepo - GetUserInfo - r.Pool.QueryRow: %w", err)
	}

	return &userInfoDTO, nil
}

// StoreUserInfo -.
func (r *UserInfoRepo) StoreUserInfo(ctx context.Context, userInfo entity.UserInfo) error {
	sql, args, err := r.Builder.
		Insert("user_credentials").
		Columns("email", "username", "password").
		Values(userInfo.Email, userInfo.Username, userInfo.Password).
		Suffix("ON CONFLICT (email, username) DO NOTHING").
		ToSql()
	if err != nil {
		return fmt.Errorf("UserInfoRepo - StoreUserInfo - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserInfoRepo - StoreUserInfo - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetUserInfo -.
func (r *UserInfoRepo) CheckUserExist(ctx context.Context, userInfo entity.UserInfo) (bool, error) {
	// Check if the user already exists
	sql, args, err := r.Builder.
		Select("1").
		From("users").
		Where(
			squirrel.Or{
				squirrel.Eq{"username": userInfo.Username},
				squirrel.Eq{"email": userInfo.Email},
			},
		).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("UserInfoRepo - CheckUserExist - r.Builder: %w", err)
	}

	var exists int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		return false, fmt.Errorf("UserInfoRepo - CheckUserExist -  r.Pool.QueryRow: %w", err)
	}

	if exists > 0 {
		return true, nil
	}

	return false, nil
}
