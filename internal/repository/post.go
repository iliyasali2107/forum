package repository

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type PostRepository interface {
	CreatePost(*models.Post) (int, error)
	GetAllPosts() ([]*models.Post, error)
	GetPostsByCategory(...int) (*[]models.Post, error)
	GetPost(int) (*models.Post, error)
	GetAllPostsOfUser(int) ([]*models.Post, error)
	UpdatePost(*models.Post) (*models.Post, error)
	DeletePost(int) error
}

type postRepo struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepo{
		db: db,
	}
}

func (r *postRepo) CreatePost(post *models.Post) (int, error) {
	query := `INSERT INTO posts(user_id, title, content, created) VALUES (?, ?, ?, ?)`
	row, err := r.db.Exec(query, post.User.ID, post.Title, post.Content, post.Created)
	if err != nil {
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	post.ID = int(id)

	return int(id), nil
}

func (r *postRepo) GetAllPosts() ([]*models.Post, error) {
	query := `SELECT * FROM posts`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	posts := []*models.Post{}
	for rows.Next() {
		post := &models.Post{User: &models.User{}}
		if err := rows.Scan(&post.ID, &post.User.ID, &post.Title, &post.Content, &post.Created); err != nil {
			fmt.Println(err)
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *postRepo) GetPostsByCategory(ids ...int) (*[]models.Post, error) {
	query := `SELECT * FROM posts WHERE id IN (?`
	for i := 0; i < len(ids); i++ {
		query += `,?`
	}
	query += `)`
	rows, err := r.db.Query(query, ids)
	if err != nil {
		return nil, err
	}
	posts := []models.Post{}
	for rows.Next() {
		post := models.Post{}
		if err := rows.Scan(&post.ID, &post.User.ID, &post.Title, &post.Content, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (r *postRepo) GetPost(id int) (*models.Post, error) {
	query := `SELECT * FROM posts WHERE id = ?`
	row := r.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	post := &models.Post{}
	if err := row.Scan(&post.ID, &post.User.ID, &post.Title, &post.Content, &post.Created); err != nil {
		return nil, err
	}

	return post, nil
}

func (r *postRepo) GetAllPostsOfUser(userID int) ([]*models.Post, error) {
	query := `SELECT * FROM posts WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	posts := []*models.Post{}
	for rows.Next() {
		post := &models.Post{User: &models.User{}}
		if err := rows.Scan(&post.ID, &post.User.ID, &post.Title, &post.Content, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *postRepo) UpdatePost(post *models.Post) (*models.Post, error) {
	return nil, nil
}

func (r *postRepo) DeletePost(id int) error {
	return nil
}
