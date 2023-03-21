package models

type Vote struct {
	ID      int
	Post    *Post
	Comment *Comment
	User    *User
	Type    string
}
