package usecase

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/pkg/validator"
	"time"
)

type createPostUsecase struct {
	postRepository     repository.PostRepository
	categoryRepository repository.CategoryRepository
}

type CreatePostUsecase interface {
	CreatePost(*validator.Validator, *models.Post) error
}

func NewCreatePostUsecase(postRepository repository.PostRepository, categoryRepository repository.CategoryRepository) CreatePostUsecase {
	return &createPostUsecase{
		postRepository:     postRepository,
		categoryRepository: categoryRepository,
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

func validatePost(v *validator.Validator, post *models.Post) {
	v.Check(post.Title != "", "title", "must be provided")
	v.Check(len(post.Title) <= 20, "title", "must not be more than 20 chars")

	v.Check(post.Content != "", "content", "must be provided")
	v.Check(len(post.Content) <= 20, "content", "must not be more than 20 chars")
}
