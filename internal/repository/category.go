package repository

import (
	"database/sql"

	"forum/internal/models"
)

type CategoryRepository interface {
	AddCategory(post_id, category_id int) error
	GetCategory(id int) (*models.Category, error)
}

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) AddCategory(post_id, category_id int) error {
	return nil
}

func (r *categoryRepo) GetCategory(id int) (*models.Category, error) {
	return nil, nil
}
