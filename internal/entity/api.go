package entity

import "time"

type RequestParams struct {
	Cursor time.Time `json:"cursor"`
	Limit  int       `json:"limit"`
	UserID string    `json:"user_uuid"`
}

type RequestParamsDTO struct {
	Cursor time.Time
	Limit  int
	UserID string
}
