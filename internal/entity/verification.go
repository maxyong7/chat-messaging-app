// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// UserCredentials -.
type UserCredentials struct {
	Username string `json:"username"  example:"username"`
	Password string `json:"password"  example:"password"`
	Email    string `json:"email"  example:"email"`
}

type UserRegistration struct {
	UserCredentials
	FirstName string `json:"first_name"  example:"first_name"`
	LastName  string `json:"last_name"  example:"last_name"`
	Avatar    string `json:"avatar,omitempty"  example:"avatar"`
}

type UserInfoDTO struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
	UserUuid string `db:"user_uuid"`
}
