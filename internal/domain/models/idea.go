package models

type Idea struct {
	ID          int64
	Image       string
	Name        string
	Description string
	Link        string
	Tags        string
	UserID      int64
	Likes       int32
	Saved       bool
}
