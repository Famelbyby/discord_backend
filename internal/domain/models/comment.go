package models

type Comment struct {
	ID           int64
	UserId       int64
	Text         string
	CreationDate string
	Avatar       string
	Username     string
}
