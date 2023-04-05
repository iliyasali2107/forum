package repository

import "database/sql"

type Repository struct {
	//AuthRepository
	CommentRepository
	PostRepository
	UserRepository
	VoteRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		CommentRepository: NewCommentRepository(db),
		PostRepository:    NewPostRepository(db),
		UserRepository:    NewUserRepository(db),
		VoteRepository:    NewVoteRepository(db),
	}
}
