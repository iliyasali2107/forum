package models

import "time"

type Post struct {
	ID            int
	User          *User
	Title         string
	Content       string
	Created       time.Time
	CreatedStr    string
	Categories    []string
	Comments      []*Comment
	Likes         int
	Dislikes      int
	CommentsCount int
}
