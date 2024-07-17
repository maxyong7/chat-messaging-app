package entity

type Reaction struct {
	UserInfo
	ReactionType string `json:"reaction_type"`
}
