package api

type PostList struct {
	Limit int     `json:"limit"`
	Total int64   `json:"total"`
	Page  int     `json:"page"`
	Items []*Post `json:"items"`
}

type CategoryList struct {
	Limit int         `json:"limit"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Items []*Category `json:"items"`
}
