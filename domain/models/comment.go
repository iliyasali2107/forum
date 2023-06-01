package models

type Comment struct {
	ID         int
	User       *User
	UserID     int
	PostID     int
	Content    string
	ParentID   int
	ReplyCount int
	Replies    []*Comment
	Likes      int
	Dislikes   int
}

