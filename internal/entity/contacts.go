package entity

type Contacts struct {
	UserInfo
	ConversationUUID string `json:"conversation_uuid"`
	Blocked          bool   `json:"blocked"`
}

type ContactsDTO struct {
	UserUUID         string `json:"user_uuid"`
	ContactUserUUID  string `json:"contact_user_uuid"`
	ConversationUUID string `json:"conversation_uuid"`
	Blocked          bool   `json:"blocked"`
	Removed          bool   `json:"removed"`
}

const DirectMessageConversationType = "direct_message"
