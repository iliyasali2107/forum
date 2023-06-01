package models

type Reaction struct {
	ID      int
	PostID  int
	CommentID int
	UserID  int
	Type    int
}

const (
	LikeType    = 1
	DislikeType = 0
)
