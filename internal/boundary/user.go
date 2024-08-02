package boundary

import "github.com/maxyong7/chat-messaging-app/internal/entity"

type RegisterUserRequestModel struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Avatar    string `json:"avatar,omitempty"`
}

type VerifyUserRequestModel struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyUserResponseModel struct {
	Token string `json:"token,omitempty"`
}

func (r RegisterUserRequestModel) ToUserRegistration() entity.UserRegistration {
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

func (r VerifyUserRequestModel) ToVerifyCredentials() entity.UserCredentials {
	return entity.UserCredentials{
		Username: r.Username,
		Password: r.Password,
	}
}
