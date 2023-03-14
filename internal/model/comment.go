package model

type Comment struct {
	ID      int
	User    *User
	Post    *Post
	Content string
	Parent  *Comment
}

