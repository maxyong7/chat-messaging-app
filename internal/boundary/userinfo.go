package boundary

import "github.com/maxyong7/chat-messaging-app/internal/entity"

type UpdateUserProfileRequestModel struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Avatar    string `json:"avatar"`
}

func (r UpdateUserProfileRequestModel) ToUserInfo(userUUID string) entity.UserInfo {
	return entity.UserInfo{
		UserUUID:  userUUID,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Avatar:    r.Avatar,
	}
}
