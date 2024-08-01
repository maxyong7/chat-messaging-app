package entity

type UserInfoDTO struct {
	UserUUID  string `json:"user_uuid,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

type UserInfo struct {
	UserUUID  string
	FirstName string
	LastName  string
	Avatar    string
}

func (dto *UserInfoDTO) ToUserInfo() UserInfo {
	if dto == nil {
		return UserInfo{}
	}
	return UserInfo{
		UserUUID:  dto.UserUUID,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Avatar:    dto.Avatar,
	}
}
