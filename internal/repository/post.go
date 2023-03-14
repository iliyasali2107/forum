package repository

import (
	"database/sql"

	"forum/internal/model"
)

type PostRepository interface {
	CreatePost(*model.Post) (int, error)
	GetAllPosts() (*[]model.Post, error)
	GetPostsByCategory(...int) (*[]model.Post, error)
	GetPost(int) (*model.Post, error)
	GetAllPostsOfUser(int) (*[]model.Post, error)
	UpdatePost(*model.Post) (*model.Post, error) // may be don't need
	DeletePost(int) error                        // may be don't need
}

type postRepo struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepo{
		db: db,
	}
}

func (r *postRepo) CreatePost(post *model.Post) (int, error) {
	query := `INSERT INTO posts(user_id, title, content, created) VALUES (?, ?, ?, ?)`
	row, err := r.db.Exec(query, post.User.ID, post.Title, post.Content, post.Created)
	if err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *postRepo) GetAllPosts() (*[]model.Post, error) {
	query := `SELECT * FROM posts`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	posts := []model.Post{}
	for rows.Next() {
		post := model.Post{}
		if err := rows.Scan(&post.ID, &post.User.ID, &post.Title, &post.Content, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (r *postRepo) GetPostsByCategory(ids ...int) (*[]model.Post, error) {
	query := `SELECT * FROM posts WHERE id IN (?`
	for i := 0; i < len(ids); i++ {
		query += `,?`
	}
	query += `)`
	rows, err := r.db.Query(query, ids)
	if err != nil {
		return nil, err
	}
	posts := []model.Post{}
	for rows.Next() {
		post := model.Post{}
		if err := rows.Scan(&post.ID, &post.User.ID, &post.Title, &post.Content, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (r *postRepo) GetPost(id int) (*model.Post, error) {
	query := `SELECT * FROM posts WHERE id = ?`
	row := r.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	post := model.Post{}
	if err := row.Scan(&post.ID, &post.User.ID, &post.Title, &post.Content, &post.Created); err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *postRepo) GetAllPostsOfUser(user_id int) (*[]model.Post, error) {
	query := `SELECT * FROM posts WHERE user_id = ?`
	rows, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}

	posts := []model.Post{}
	for rows.Next() {
		post := model.Post{}
		if err := rows.Scan(&post.ID, &post.User.ID, &post.Title, &post.Content, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (r *postRepo) UpdatePost(post *model.Post) (*model.Post, error) {
	return nil, nil
}

func (r *postRepo) DeletePost(id int) error {
	return nil
}
