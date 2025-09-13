package models

type Message struct {
	ID           int64
	SenderId     int64
	RecieverId   int64
	Filename     string
	Text         string
	CreationDate string
	IdeaId       int64
}
