package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
)

type CommentRepository interface {
	CreateCommentWithParent(*models.Comment) (int, error)
	CreateCommentWithoutParent(*models.Comment) (int, error)
	GetComment(int) (*models.Comment, error)
	GetPostComments(int) (*[]models.Comment, error)
	GetUserComments(int) (*[]models.Comment, error)
	GetAllComments() (*[]models.Comment, error)             // may be don't need
	UpdateComment(*models.Comment) (*models.Comment, error) // may be don't need
	DeleteComment(int) error                                // may be don't need
}

type commentRepo struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepo{
		db: db,
	}
}

func (r *commentRepo) CreateCommentWithParent(comment *models.Comment) (int, error) {
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

func (r *commentRepo) CreateCommentWithoutParent(comment *models.Comment) (int, error) {
	fmt.Println(comment)
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

func (r *commentRepo) GetComment(id int) (*models.Comment, error) {
	query := `SELECT * FROM comments WHERE id = ?`
	row := r.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	comment := models.Comment{}
	if err := row.Scan(&comment.ID, &comment.User.ID, &comment.Post.ID, &comment.Parent.ID); err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepo) GetAllComments() (*[]models.Comment, error) {
	query := `SELECT * FROM comments`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	comments := []models.Comment{}
	for rows.Next() {
		comment := models.Comment{}
		if err := rows.Scan(&comment.ID, &comment.User.ID, &comment.Post.ID, &comment.Content, &comment.Parent.ID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return &comments, nil
}

func (r *commentRepo) GetPostComments(postID int) (*[]models.Comment, error) {
	query := `SELECT * FROM comments WHERE post_id = ?`
	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}

	comments := []models.Comment{}
	for rows.Next() {
		comment := models.Comment{}
		if err := rows.Scan(&comment.ID, &comment.User.ID, &comment.Post.ID, &comment.Content, &comment.Parent.ID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return &comments, nil
}

func (r *commentRepo) GetUserComments(userID int) (*[]models.Comment, error) {
	query := `SELECT * FROM comments WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	comments := []models.Comment{}
	for rows.Next() {
		comment := models.Comment{}
		if err := rows.Scan(&comment.ID, &comment.User.ID, &comment.Post.ID, &comment.Content, &comment.Parent.ID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return &comments, nil
}

func (r *commentRepo) UpdateComment(*models.Comment) (*models.Comment, error) {
	return nil, nil
}

func (r *commentRepo) DeleteComment(id int) error {
	return nil
}
