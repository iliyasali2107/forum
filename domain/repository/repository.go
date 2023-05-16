package repository

import "database/sql"

type Repository struct {
	// AuthRepository
	CommentRepository
	PostRepository
	UserRepository
	ReactionRepository
	CategoryRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		CommentRepository:  NewCommentRepository(db),
		PostRepository:     NewPostRepository(db),
		UserRepository:     NewUserRepository(db),
		ReactionRepository: NewReactionRepository(db),
		CategoryRepository: NewCategoryRepository(db),
	}
}
