package repository

import (
	"database/sql"

	"forum/domain/models"
)

type CategoryRepository interface {
	AddCategory(postID, categoryID int) error
	GetCategory(name string) (*models.Category, error)
	GetAllCategories() ([]*models.Category, error)
	GetCategoriesForPost(post *models.Post) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (cr *categoryRepository) AddCategory(postID, categoryID int) error {
	query := `INSERT INTO categories_posts(post_id, category_id) VALUES (?, ?)`
	row, err := cr.db.Exec(query, postID, categoryID)
	if err != nil {
		return err
	}

	_, err = row.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (cr *categoryRepository) GetCategory(name string) (*models.Category, error) {
	query := `SELECT * FROM categories WHERE name = ?`
	row := cr.db.QueryRow(query, name)
	if row.Err() != nil {
		return nil, row.Err()
	}

	category := &models.Category{}

	if err := row.Scan(&category.ID, &category.Name); err != nil {
		return nil, err
	}

	return category, nil
}

func (cr *categoryRepository) GetAllCategories() ([]*models.Category, error) {
	query := `SELECT * FROM categories`
	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, err
	}

	categories := []*models.Category{}
	for rows.Next() {
		category := &models.Category{}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)

	}

	return categories, nil
}

func (cr *categoryRepository) GetCategoriesForPost(post *models.Post) error {
	query := `SELECT name
			FROM categories
			INNER JOIN categories_posts
			ON categories.id = categories_posts.category_id
			WHERE post_id = ?;`

	rows, err := cr.db.Query(query, post.ID)
	if err != nil {
		return err
	}
	for rows.Next() {
		category := ""

		if err = rows.Scan(&category); err != nil {
			return err
		}

		post.Categories = append(post.Categories, category)
	}
	return nil
}
