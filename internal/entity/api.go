package entity

import "time"

// UserInfo -.
type Pagination struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

type RequestParams struct {
	Cursor time.Time `json:"cursor"`
	Limit  int       `json:"limit"`
	UserID string    `json:"user_uuid"`
}
