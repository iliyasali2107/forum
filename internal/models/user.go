package models

import "time"

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Token    string
	Expires  time.Time
}

var AnonymousUser = &User{}
