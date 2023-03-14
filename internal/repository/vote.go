package repository

import (
	"database/sql"

	"forum/internal/model"
)

type VoteRepository interface {
	CreateCommentVote(*model.Vote) (int, error)
	CreatePostVote(*model.Vote) (int, error)
	GetPostUpvotes(int) (int, error)
	GetPostDownvotes(int) (int, error)
	GetCommentUpvotes(int) (int, error)
	GetCommentDownvotes(int) (int, error)
	DeleteCommentVote(int) error
	DeletePostVote(int) error
}

type voteRepo struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) VoteRepository {
	return &voteRepo{
		db: db,
	}
}

func (r *voteRepo) CreateCommentVote(vote *model.Vote) (int, error) {
	query := `INSERT INTO votes_comments (comment_id, user_id, type) VALUES (?, ?, ?)`
	row, err := r.db.Exec(query, vote.Comment.ID, vote.User.ID, vote.Type)
	if err != nil {
		return 0, nil
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (r *voteRepo) CreatePostVote(vote *model.Vote) (int, error) {
	query := `INSERT INTO votes_posts (post_id, user_id, type) VALUES (?, ?, ?)`
	row, err := r.db.Exec(query, vote.Post.ID, vote.User.ID, vote.Type)
	if err != nil {
		return 0, nil
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (r *voteRepo) GetPostUpvotes(post_id int) (int, error) {
	query := `SELECT count() FROM votes_posts WHERE post_id = ? AND type = 1`
	var upvotes int
	row := r.db.QueryRow(query, post_id)
	if err := row.Scan(&upvotes); err != nil {
		return 0, err
	}

	return upvotes, nil
}

func (r *voteRepo) GetPostDownvotes(post_id int) (int, error) {
	query := `SELECT count() FROM votes_posts WHERE post_id = ? AND type = 0`
	var downvotes int
	row := r.db.QueryRow(query, post_id)
	if err := row.Scan(&downvotes); err != nil {
		return 0, err
	}

	return downvotes, nil
}

func (r *voteRepo) GetCommentUpvotes(comment_id int) (int, error) {
	query := `SELECT count() FROM votes_comments WHERE comment_id = ? AND type = 1`
	var upvotes int
	row := r.db.QueryRow(query, comment_id)
	if err := row.Scan(&upvotes); err != nil {
		return 0, err
	}

	return upvotes, nil
}

func (r *voteRepo) GetCommentDownvotes(comment_id int) (int, error) {
	query := `SELECT count() FROM votes_comments WHERE comment_id = ? AND type = 0`
	var downvotes int
	row := r.db.QueryRow(query, comment_id)
	if err := row.Scan(&downvotes); err != nil {
		return 0, err
	}

	return downvotes, nil
}

func (r *voteRepo) DeleteCommentVote(comment_id int) error {
	query := `DELETE FROM votes_comments WHERE comment_id = ?`
	if _, err := r.db.Exec(query, comment_id); err != nil {
		return err
	}
	return nil
}

func (r *voteRepo) DeletePostVote(post_id int) error {
	query := `DELETE FROM votes_comments WHERE post_id = ?`
	if _, err := r.db.Exec(query, post_id); err != nil {
		return err
	}
	return nil
}
