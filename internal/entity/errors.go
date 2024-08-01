package entity

import "errors"

var (
	ErrUserAlreadyExists          = errors.New("user already exists")
	ErrContactAlreadyExists       = errors.New("contact already exists")
	ErrContactDoesNotExists       = errors.New("contact does not exists")
	ErrNoConversationFound        = errors.New("no conversations found")
	ErrUserNameNotFound           = errors.New("no username found")
	ErrUserNotFound               = errors.New("no user found")
	ErrUserNotInGroupChat         = errors.New("user not in groupchat")
	ErrParticipantAlrdInGroupChat = errors.New("one or more participant(s) is already in groupchat")
	ErrParticipantNotInGroupChat  = errors.New("one or more participant(s) is not in groupchat")
	ErrIncorrectPassword          = errors.New("incorrect password")

	// Add other custom errors here
)
