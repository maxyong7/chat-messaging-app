// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// UserCredentials -.
type UserCredentials struct {
	Username string
	Password string
	Email    string
}

type UserRegistration struct {
	UserCredentials
	FirstName string
	LastName  string
	Avatar    string
}

type UserCredentialsDTO struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
	UserUuid string `db:"user_uuid"`
}
