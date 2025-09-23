package models

type User struct {
	ID           uint32
	Email        string
	PasswordHash []byte
}
