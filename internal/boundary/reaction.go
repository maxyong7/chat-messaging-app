package boundary

type MessageReactionMenu struct {
	MessageUUID  string `json:"message_uuid"`
	ReactionType string `json:"reaction_type"`
}
type RemoveReactionRequest struct {
	MessageUUID string `json:"message_uuid"`
}

type ReactionResponseData struct {
	MessageUUID string `json:"message_uuid"`
	UserUUID    string `json:"user_uuid"`
	Reaction    string `json:"reaction,omitempty"`
}
