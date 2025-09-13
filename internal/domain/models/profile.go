package models

type Profile struct {
	ID          int64
	Email       string
	AvatarImage string
	Name        string
	Description string
	Link        string
	Boards      []int64
	SavedIdeas  []*IdeaBoardPair
}
type ProfileLight struct {
	ID          int64
	Name        string
	AvatarImage string
}
type IdeaBoardPair struct {
	IdeaId  int64
	BoardId int64
}
