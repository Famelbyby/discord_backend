package models

type Board struct {
	ID       int64
	Name     string
	IdeasIds []int64
}