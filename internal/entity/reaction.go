package entity

type Reaction struct {
	MessageUUID  string
	SenderUUID   string
	ReactionType string
}
type GetReactionDTO struct {
	UserProfileDTO
	ReactionType string `json:"reaction_type"`
}

type StoreReactionDTO struct {
	MessageUUID  string
	SenderUUID   string
	ReactionType string
}

type RemoveReactionDTO struct {
	MessageUUID string
	SenderUUID  string
}
