package boundary

import "github.com/maxyong7/chat-messaging-app/internal/entity"

type RegistrationForm struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Avatar    string `json:"avatar,omitempty"`
}

type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LogoutScreen struct {
	Token string `json:"token,omitempty"`
}

func (r RegistrationForm) ToUserRegistration() entity.UserRegistration {
	return entity.UserRegistration{
		UserCredentials: entity.UserCredentials{
			Username: r.Username,
			Password: r.Password,
			Email:    r.Email,
		},
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Avatar:    r.Avatar,
	}
}

func (r LoginForm) ToUserCredentials() entity.UserCredentials {
	return entity.UserCredentials{
		Username: r.Username,
		Password: r.Password,
	}
}
