package models

type Reaction struct {
	ID      int
	Post    *Post
	Comment *Comment
	User    *User
	Type    int
}

const (
	LikeType    = 1
	DislikeType = 0
)
