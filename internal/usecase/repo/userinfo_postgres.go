package repo

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
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

// GetUserCredentials -.
func (r *UserInfoRepo) GetUserCredentials(ctx context.Context, userInfo entity.UserCredentials) (*entity.UserCredentialsDTO, error) {
	sql := `
		SELECT email, username, password, user_uuid
		FROM user_credentials
		WHERE (username = $1 OR email = $2) 
	`

	var userInfoDTO entity.UserCredentialsDTO
	err := r.Pool.QueryRow(ctx, sql, userInfo.Username, userInfo.Email).
		Scan(&userInfoDTO.Email, &userInfoDTO.Username, &userInfoDTO.Password, &userInfoDTO.UserUuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("UserInfoRepo - GetUserCredentials - r.Pool.QueryRow: %w", err)
	}

	return &userInfoDTO, nil
}

// StoreUserInfo -.
func (r *UserInfoRepo) StoreUserInfo(ctx context.Context, userRegis entity.UserRegistration) error {
	userUuid := uuid.New()
	// Begin a transaction
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("UserInfoRepo - StoreUserInfo - failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back if it doesn't commit
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			tx.Rollback(ctx) // err is non-nil; rollback
		}
	}()

	insertUserCredentialsSQL := `
		INSERT INTO user_credentials (email, username, password, user_uuid)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, insertUserCredentialsSQL, userRegis.Email, userRegis.Username, userRegis.Password, userUuid)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertUserCredentialsSQL query: %w", err)
	}

	insertUserInfoSQL := `
		INSERT INTO user_info (user_uuid, first_name, last_name, email, avatar)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, insertUserInfoSQL, userUuid, userRegis.FirstName, userRegis.LastName, userRegis.Email, userRegis.Avatar)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertUserInfoSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("UserInfoRepo - StoreUserInfo - failed to commit transaction: %w", err)
	}

	return nil
}

// CheckUserExist -.
func (r *UserInfoRepo) CheckUserExist(ctx context.Context, userRegis entity.UserRegistration) (bool, error) {
	// Check if the user already exists
	sql, args, err := r.Builder.
		Select("1").
		From("user_credentials").
		Where(
			squirrel.Or{
				squirrel.Eq{"username": userRegis.Username},
				squirrel.Eq{"email": userRegis.Email},
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

// GetUserInfo -.
func (r *UserInfoRepo) GetUserInfo(ctx context.Context, userUuid string) (*entity.UserInfoDTO, error) {
	sql := `
		SELECT user_uuid, first_name, last_name, avatar
		FROM user_info
		WHERE (user_uuid = $1) 
	`

	var userInfoDTO entity.UserInfoDTO
	err := r.Pool.QueryRow(ctx, sql, userUuid).
		Scan(&userInfoDTO.UserUUID, &userInfoDTO.FirstName, &userInfoDTO.LastName, &userInfoDTO.Avatar)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("UserInfoRepo - GetUserCredentials - r.Pool.QueryRow: %w", err)
	}

	return &userInfoDTO, nil
}

// GetUserUUIDByUsername -.
func (r *UserInfoRepo) GetUserUUIDByUsername(ctx context.Context, userName string) (*string, error) {
	sql := `
		SELECT user_uuid
		FROM user_credentials
		WHERE (username = $1) 
	`

	var userUUID string
	err := r.Pool.QueryRow(ctx, sql, userName).
		Scan(&userUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("UserInfoRepo - GetUserCredentials - r.Pool.QueryRow: %w", err)
	}

	return &userUUID, nil
}
