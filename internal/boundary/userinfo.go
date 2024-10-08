package boundary

import "github.com/maxyong7/chat-messaging-app/internal/entity"

type ProfileSettingScreen struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Avatar    string `json:"avatar"`
}

func (r ProfileSettingScreen) ToUserInfo(userUUID string) entity.UserProfile {
	return entity.UserProfile{
		UserUUID:  userUUID,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Avatar:    r.Avatar,
	}
}
