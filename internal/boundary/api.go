package boundary

type Pagination struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}
