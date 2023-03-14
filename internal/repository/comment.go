package repository

import (
	"database/sql"

	"forum/internal/model"
)

type CommentRepository interface {
	CreateCommentWithParent(*model.Comment) (int, error)
	CreateCommentWithoutParent(*model.Comment) (int, error)
	GetComment(int) (*model.Comment, error)
	GetPostComments(int) (*[]model.Comment, error)
	GetUserComments(int) (*[]model.Comment, error)
	GetAllComments() (*[]model.Comment, error)            // may be don't need
	UpdateComment(*model.Comment) (*model.Comment, error) // may be don't need
	DeleteComment(int) error                              // may be don't need
}

type commentRepo struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepo{
		db: db,
	}
}

func (r *commentRepo) CreateCommentWithParent(comment *model.Comment) (int, error) {
	query := `INSERT INTO comments (user_id, post_id, content, parent_id) VALUES (?, ?, ?, ?)`
	row, err := r.db.Exec(query, comment.User.ID, comment.Post.ID, comment.Content, comment.Parent.ID)
	if err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *commentRepo) CreateCommentWithoutParent(comment *model.Comment) (int, error) {
	query := `INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)`
	row, err := r.db.Exec(query, comment.User.ID, comment.Post.ID, comment.Content)
	if err != nil {
		return 0, err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *commentRepo) GetComment(id int) (*model.Comment, error) {
	query := `SELECT * FROM comments WHERE id = ?`
	row := r.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	comment := model.Comment{}
	if err := row.Scan(&comment.ID, &comment.User.ID, &comment.Post.ID, &comment.Parent.ID); err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepo) GetAllComments() (*[]model.Comment, error) {
	query := `SELECT * FROM comments`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	comments := []model.Comment{}
	for rows.Next() {
		comment := model.Comment{}
		if err := rows.Scan(&comment.ID, &comment.User.ID, &comment.Post.ID, &comment.Content, &comment.Parent.ID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return &comments, nil
}

func (r *commentRepo) GetPostComments(post_id int) (*[]model.Comment, error) {
	query := `SELECT * FROM comments WHERE post_id = ?`
	rows, err := r.db.Query(query, post_id)
	if err != nil {
		return nil, err
	}

	comments := []model.Comment{}
	for rows.Next() {
		comment := model.Comment{}
		if err := rows.Scan(&comment.ID, &comment.User.ID, &comment.Post.ID, &comment.Content, &comment.Parent.ID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return &comments, nil
}

func (r *commentRepo) GetUserComments(user_id int) (*[]model.Comment, error) {
	query := `SELECT * FROM comments WHERE user_id = ?`
	rows, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}

	comments := []model.Comment{}
	for rows.Next() {
		comment := model.Comment{}
		if err := rows.Scan(&comment.ID, &comment.User.ID, &comment.Post.ID, &comment.Content, &comment.Parent.ID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return &comments, nil
}

func (r *commentRepo) UpdateComment(*model.Comment) (*model.Comment, error) {
	return nil, nil
}

func (r *commentRepo) DeleteComment(id int) error {
	return nil
}
