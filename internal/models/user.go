package models

import "time"

type User struct {
	ID       int
	Name     string
	Email    string
	Password Password
	Token    string
	Expires  time.Time
}

type Password struct {
	Plaintext *string
	Hash      *[]byte
}
