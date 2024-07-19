package entity

type GroupChatRequest struct {
	UserUUID         string        `json:"user_uuid"`
	Title            string        `json:"title"`
	ConversationUUID string        `json:"conversation_uuid"`
	Participants     []Participant `json:"participants"`
}

type Participant struct {
	ParticipantUUID string `json:"participant_uuid"`
}

const GroupMessageConversationType = "group_message"
