package entity

type UserProfileDTO struct {
	UserUUID  string `json:"user_uuid,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

type UserProfile struct {
	UserUUID  string
	FirstName string
	LastName  string
	Avatar    string
}

func (dto *UserProfileDTO) ToUserInfo() UserProfile {
	if dto == nil {
		return UserProfile{}
	}
	return UserProfile{
		UserUUID:  dto.UserUUID,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Avatar:    dto.Avatar,
	}
}
