package repository

import (
	"database/sql"

	"forum/domain/models"
)

type CommentRepository interface {
	CreateComment(*models.Comment) (int, error)
	GetComment(int) (*models.Comment, error)
	GetPostComments(int) ([]*models.Comment, error)
	GetCommentReplies(int) ([]*models.Comment, error)
	GetCommentRepliesCount(int) (int, error)
	GetUserComments(int) (*[]models.Comment, error)
	GetAllComments() (*[]models.Comment, error)             // may be don't need
	UpdateComment(*models.Comment) (*models.Comment, error) // may be don't need
	DeleteComment(int) error                                // may be don't need
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (cr *commentRepository) CreateComment(comment *models.Comment) (int, error) {
	query := `INSERT INTO comments (user_id, post_id, content, parent_id) VALUES (?, ?, ?, ?)`
	row, err := cr.db.Exec(query, comment.UserID, comment.PostID, comment.Content, comment.ParentID)
	if err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (cr *commentRepository) GetComment(id int) (*models.Comment, error) {
	query := `SELECT * FROM comments WHERE id = ?`
	row := cr.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	comment := models.Comment{}
	if err := row.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.ParentID); err != nil {
		return nil, err
	}
	return &comment, nil
}

func (cr *commentRepository) GetAllComments() (*[]models.Comment, error) {
	query := `SELECT * FROM comments`
	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, err
	}

	comments := []models.Comment{}
	for rows.Next() {
		comment := models.Comment{}
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.ParentID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return &comments, nil
}

func (cr *commentRepository) GetPostComments(postID int) ([]*models.Comment, error) {
	query := `SELECT * FROM comments WHERE post_id = ? AND parent_id = 0`
	rows, err := cr.db.Query(query, postID)
	if err != nil {
		return nil, err
	}

	comments := []*models.Comment{}
	for rows.Next() {
		comment := &models.Comment{}
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.ParentID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (cr *commentRepository) GetCommentRepliesCount(parentID int) (int, error) {
	query := `SELECT count() FROM comments WHERE parent_id = ?`
	var replies int
	row := cr.db.QueryRow(query, parentID)
	if err := row.Scan(&replies); err != nil {
		return 0, err
	}

	return replies, nil
}

func (cr *commentRepository) GetCommentReplies(parentID int) ([]*models.Comment, error) {
	query := `SELECT * FROM comments WHERE parent_id = ?`
	rows, err := cr.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}

	comments := []*models.Comment{}
	for rows.Next() {
		comment := &models.Comment{}
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.ParentID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (cr *commentRepository) GetUserComments(userID int) (*[]models.Comment, error) {
	query := `SELECT * FROM comments WHERE user_id = ?`
	rows, err := cr.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	comments := []models.Comment{}
	for rows.Next() {
		comment := models.Comment{}
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.ParentID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return &comments, nil
}

func (cr *commentRepository) UpdateComment(*models.Comment) (*models.Comment, error) {
	return nil, nil
}

func (cr *commentRepository) DeleteComment(id int) error {
	return nil
}
