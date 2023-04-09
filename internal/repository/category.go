package repository

import (
	"database/sql"

	"forum/internal/models"
)

type CategoryRepository interface {
	AddCategory(postID, categoryID int) error
	GetCategory(name string) (*models.Category, error)
	GetAllCategories() ([]*models.Category, error)
}

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) AddCategory(postID, categoryID int) error {
	query := `INSERT INTO categories_posts(post_id, category_id) VALUES (?, ?)`
	row, err := r.db.Exec(query, postID, categoryID)
	if err != nil {
		return err
	}

	_, err = row.LastInsertId()
	if err != nil {
		return err
	}

	return nil

}

func (r *categoryRepo) GetCategory(name string) (*models.Category, error) {
	query := `SELECT * FROM categories WHERE name = ?`
	row := r.db.QueryRow(query, name)
	if row.Err() != nil {
		return nil, row.Err()
	}

	category := &models.Category{}

	if err := row.Scan(&category.ID, &category.Name); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *categoryRepo) GetAllCategories() ([]*models.Category, error) {
	query := `SELECT * FROM categories`
	rows, err := r.db.Query(query)
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
