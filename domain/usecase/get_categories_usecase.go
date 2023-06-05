package usecase

import (
	"fmt"
	"forum/domain/models"
	"forum/domain/repository"
)

type getCategoriesUsecase struct {
	categoryRepository repository.CategoryRepository
}

type GetCategoriesUsecase interface {
	GetAllCategories() ([]*models.Category, error)
}

func NewGetCategoriesUsecase(categoryRepository repository.CategoryRepository) GetCategoriesUsecase {
	return &getCategoriesUsecase{
		categoryRepository: categoryRepository,
	}
}

func (gcu *getCategoriesUsecase) GetAllCategories() ([]*models.Category, error) {
	categories, err := gcu.categoryRepository.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("couldn't get all catgories: %w", err)
	}

	return categories, nil
}
