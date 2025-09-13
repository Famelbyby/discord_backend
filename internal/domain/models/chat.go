package models

type ChatData struct {
	ID         int64
	FirstData  ChatUserData
	SecondData ChatUserData
}
type ChatUserData struct {
	UserId   int64
	Username string
	Avatar   string
}
