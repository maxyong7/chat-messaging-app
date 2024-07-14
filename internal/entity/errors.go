package entity

import "errors"

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrContactAlreadyExists = errors.New("contact already exists")
	ErrNoConversationFound  = errors.New("no conversations found")
	ErrUserNameNotFound     = errors.New("no username found")

	// Add other custom errors here
)
