package model

import "time"

type Post struct {
	ID      int
	User    *User
	Title   string
	Content string
	Created time.Time
	Categories *[]Category
}

