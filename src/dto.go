package api

type PostList struct {
	Limit int
	Total int
	Page int
	Items []*Post
}