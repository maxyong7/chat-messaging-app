package entity

import "errors"

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrNoConversationFound = errors.New("no conversations found")

	// Add other custom errors here
)
