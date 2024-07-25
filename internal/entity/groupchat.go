package entity

type GroupChat struct {
	UserUUID         string        `json:"user_uuid"`
	Title            string        `json:"title"`
	ConversationUUID string        `json:"conversation_uuid"`
	Participants     []Participant `json:"participants"`
}

type Participant struct {
	Username        string `json:"username"`
	ParticipantUUID string `json:"participant_uuid"`
}

type GroupChatDTO struct {
	UserUUID         string
	Title            string
	ConversationUUID string
	Participants     []ParticipantDTO
}

type ParticipantDTO struct {
	Username        string
	ParticipantUUID string
}

const GroupMessageConversationType = "group_message"
