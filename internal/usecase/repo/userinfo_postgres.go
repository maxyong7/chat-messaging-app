package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/maxyong7/chat-messaging-app/internal/entity"
)

// UserInfoRepo -.
type UserInfoRepo struct {
	// *postgres.Postgres
	*sql.DB
}

// New -.
func NewUserInfo(pg *sql.DB) *UserInfoRepo {
	return &UserInfoRepo{pg}
}

// GetUserCredentials -.
func (r *UserInfoRepo) GetUserCredentials(ctx context.Context, userInfo entity.UserCredentials) (*entity.UserCredentialsDTO, error) {
	getUserCredentialsSQL := `
		SELECT email, username, password, user_uuid
		FROM user_credentials
		WHERE (username = $1 OR email = $2) 
	`

	var userInfoDTO entity.UserCredentialsDTO
	err := r.QueryRowContext(ctx, getUserCredentialsSQL, userInfo.Username, userInfo.Email).
		Scan(&userInfoDTO.Email, &userInfoDTO.Username, &userInfoDTO.Password, &userInfoDTO.UserUuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("UserInfoRepo - GetUserCredentials - r.QueryRowContext: %w", err)
	}

	return &userInfoDTO, nil
}

// StoreUserInfo -.
func (r *UserInfoRepo) StoreUserInfo(ctx context.Context, userRegis entity.UserRegistration) error {
	userUuid := uuid.New()
	// Begin a transaction
	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("UserInfoRepo - StoreUserInfo - failed to begin transaction: %w", err)
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

	insertUserCredentialsSQL := `
		INSERT INTO user_credentials (email, username, password, user_uuid)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.ExecContext(ctx, insertUserCredentialsSQL, userRegis.Email, userRegis.Username, userRegis.Password, userUuid)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertUserCredentialsSQL query: %w", err)
	}

	insertUserInfoSQL := `
		INSERT INTO user_info (user_uuid, first_name, last_name, email, avatar)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.ExecContext(ctx, insertUserInfoSQL, userUuid, userRegis.FirstName, userRegis.LastName, userRegis.Email, userRegis.Avatar)
	if err != nil {
		return fmt.Errorf("failed to execute insert insertUserInfoSQL query: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("UserInfoRepo - StoreUserInfo - failed to commit transaction: %w", err)
	}

	return nil
}

// CheckUserExist -.
func (r *UserInfoRepo) CheckUserExist(ctx context.Context, userRegis entity.UserRegistration) (bool, error) {
	// Check if the user already exists
	checkUserExistSQL := `
	SELECT 1 
	FROM user_credentials
	WHERE (username = $1 or email = $2)
	`
	// sql, args, err := r.Builder.
	// 	Select("1").
	// 	From("user_credentials").
	// 	Where(
	// 		squirrel.Or{
	// 			squirrel.Eq{"username": userRegis.Username},
	// 			squirrel.Eq{"email": userRegis.Email},
	// 		},
	// 	).
	// 	ToSql()

	// if err != nil {
	// 	return false, fmt.Errorf("UserInfoRepo - CheckUserExist - r.Builder: %w", err)
	// }

	var exists int
	err := r.QueryRowContext(ctx, checkUserExistSQL, userRegis.Username, userRegis.Email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("UserInfoRepo - CheckUserExist -  r.QueryRowContext: %w", err)
	}

	if exists > 0 {
		return true, nil
	}

	return false, nil
}

// GetUserInfo -.
func (r *UserInfoRepo) GetUserInfo(ctx context.Context, userUuid string) (*entity.UserInfoDTO, error) {
	getUserInfoSQL := `
		SELECT user_uuid, first_name, last_name, avatar
		FROM user_info
		WHERE (user_uuid = $1) 
	`

	var userInfoDTO entity.UserInfoDTO
	err := r.QueryRowContext(ctx, getUserInfoSQL, userUuid).
		Scan(&userInfoDTO.UserUUID, &userInfoDTO.FirstName, &userInfoDTO.LastName, &userInfoDTO.Avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("UserInfoRepo - GetUserCredentials - r.QueryRowContext: %w", err)
	}

	return &userInfoDTO, nil
}

// GetUserUUIDByUsername -.
func (r *UserInfoRepo) GetUserUUIDByUsername(ctx context.Context, userName string) (*string, error) {
	getUserUUIDByUsernameSQL := `
		SELECT user_uuid
		FROM user_credentials
		WHERE (username = $1) 
	`

	var userUUID string
	err := r.QueryRowContext(ctx, getUserUUIDByUsernameSQL, userName).
		Scan(&userUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("UserInfoRepo - GetUserCredentials - r.QueryRowContext: %w", err)
	}

	return &userUUID, nil
}
