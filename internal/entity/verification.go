// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Verification -.
type Verification struct {
	Username string `json:"username"  example:"username"`
	Password string `json:"password"  example:"password"`
	Email    string `json:"email"  example:"email"`
}

type UserInfoDTO struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}
