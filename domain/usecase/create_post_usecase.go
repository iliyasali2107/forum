package usecase

import (
	"time"

	"forum/domain/models"
	"forum/domain/repository"
	"forum/pkg/validator"
)

type createPostUsecase struct {
	postRepository     repository.PostRepository
	categoryRepository repository.CategoryRepository
	ContextTimeout     time.Duration
}

type CreatePostUsecase interface {
	CreatePost(*validator.Validator, *models.Post) error
	GetAllCategories() ([]*models.Category, error)
}

func NewCreatePostUsecase(postRepository repository.PostRepository, categoryRepository repository.CategoryRepository, timeout time.Duration) CreatePostUsecase {
	return &createPostUsecase{
		postRepository:     postRepository,
		categoryRepository: categoryRepository,
		ContextTimeout:     timeout,
	}
}

func (cpu *createPostUsecase) CreatePost(v *validator.Validator, post *models.Post) error {
	post.Created = time.Now()
	if validatePost(v, post); !v.Valid() {
		return ErrFormValidation
	}

	_, err := cpu.postRepository.CreatePost(post)
	if err != nil {
		return err
	}

	for _, c := range post.Categories {
		category, err := cpu.categoryRepository.GetCategory(c)
		if err != nil {
			return ErrCategoryNotFound
		}

		err = cpu.categoryRepository.AddCategory(post.ID, category.ID)
		if err != nil {
			return err
		}
	}
	return err
}

func (cpu *createPostUsecase) GetAllCategories() ([]*models.Category, error) {
	categories, err := cpu.categoryRepository.GetAllCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func validatePost(v *validator.Validator, post *models.Post) {
	v.Check(post.Title != "", "title", "must be provided")
	v.Check(len(post.Title) <= 20, "title", "must not be more than 20 chars")

	v.Check(post.Content != "", "content", "must be provided")
	v.Check(len(post.Content) <= 20, "content", "must not be more than 20 chars")
}
