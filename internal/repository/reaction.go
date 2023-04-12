package repository

import (
	"database/sql"

	"forum/internal/models"
)

type ReactionRepository interface {
	CreateCommentVote(reaction *models.Reaction) (int, error)
	CreatePostReaction(reaction *models.Reaction) error
	GetPostLikes(int) (int, error)
	GetPostDislikes(int) (int, error)
	GetCommentLikes(int) (int, error)
	GetCommentDislikes(int) (int, error)
	DeleteCommentReaction(int) error
	DeletePostReaction(reaction *models.Reaction) error
	GetReactionOfPost(reaction *models.Reaction) (*models.Reaction, error)
}

type reactionRepo struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) ReactionRepository {
	return &reactionRepo{
		db: db,
	}
}

func (r *reactionRepo) CreateCommentVote(reaction *models.Reaction) (int, error) {
	query := `INSERT INTO reactions_comments (comment_id, user_id, type) VALUES (?, ?, ?)`
	row, err := r.db.Exec(query, reaction.Comment.ID, reaction.User.ID, reaction.Type)
	if err != nil {
		return 0, nil
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (r *reactionRepo) CreatePostReaction(reaction *models.Reaction) error {
	query := `INSERT INTO reactions_posts (post_id, user_id, type) VALUES (?, ?, ?)`
	row, err := r.db.Exec(query, reaction.Post.ID, reaction.User.ID, reaction.Type)
	if err != nil {
		return nil
	}

	_, err = row.LastInsertId()
	if err != nil {
		return err
	}

	return err
}

func (r *reactionRepo) GetPostLikes(post_id int) (int, error) {
	query := `SELECT count() FROM reactions_posts WHERE post_id = ? AND type = 1`
	var likes int
	row := r.db.QueryRow(query, post_id)
	if err := row.Scan(&likes); err != nil {
		return 0, err
	}

	return likes, nil
}

func (r *reactionRepo) GetPostDislikes(post_id int) (int, error) {
	query := `SELECT count() FROM reactions_posts WHERE post_id = ? AND type = 0`
	var dislikes int
	row := r.db.QueryRow(query, post_id)
	if err := row.Scan(&dislikes); err != nil {
		return 0, err
	}

	return dislikes, nil
}

func (r *reactionRepo) GetCommentLikes(comment_id int) (int, error) {
	query := `SELECT count() FROM reactions_comments WHERE comment_id = ? AND type = 1`
	var likes int
	row := r.db.QueryRow(query, comment_id)
	if err := row.Scan(&likes); err != nil {
		return 0, err
	}

	return likes, nil
}

func (r *reactionRepo) GetCommentDislikes(comment_id int) (int, error) {
	query := `SELECT count() FROM reactions_comments WHERE comment_id = ? AND type = 0`
	var dislikes int
	row := r.db.QueryRow(query, comment_id)
	if err := row.Scan(&dislikes); err != nil {
		return 0, err
	}

	return dislikes, nil
}

func (r *reactionRepo) DeleteCommentReaction(comment_id int) error {
	query := `DELETE FROM reactions_comments WHERE comment_id = ?`
	if _, err := r.db.Exec(query, comment_id); err != nil {
		return err
	}
	return nil
}

func (r *reactionRepo) DeletePostReaction(reaction *models.Reaction) error {
	query := `DELETE FROM reactions_comments WHERE post_id = ? AND user_id = ?`
	if _, err := r.db.Exec(query, reaction.Post.ID, reaction.User.ID); err != nil {
		return err
	}
	return nil
}

func (r *reactionRepo) GetReactionOfPost(reaction *models.Reaction) (*models.Reaction, error) {
	query := `SELECT * FROM reactions_posts WHERE user_id = ? AND post_id = ?`

	row := r.db.QueryRow(query, reaction.User.ID, reaction.Post.ID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	react := &models.Reaction{}
	if err := row.Scan(&react.ID, &react.Post.ID, &react.User.ID, &react.Type); err != nil {
		return nil, err
	}

	return react, nil
}
