package entity

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

type UserRegistrationDTO struct {
	Username  string
	Password  string
	Email     string
	FirstName string
	LastName  string
	Avatar    string
}

type UserCredentialsDTO struct {
	ID       int
	Username string
	Password string
	Email    string
	UserUuid string
}
