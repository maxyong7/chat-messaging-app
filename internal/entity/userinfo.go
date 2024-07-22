package entity

type UserInfoDTO struct {
	UserUUID  string
	FirstName string
	LastName  string
	Avatar    string
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
