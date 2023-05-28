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
	CreatePost(*validator.Validator, *models.Post) (int, error)
	GetAllCategories() ([]*models.Category, error)
}

func NewCreatePostUsecase(postRepository repository.PostRepository, categoryRepository repository.CategoryRepository, timeout time.Duration) CreatePostUsecase {
	return &createPostUsecase{
		postRepository:     postRepository,
		categoryRepository: categoryRepository,
		ContextTimeout:     timeout,
	}
}

func (cpu *createPostUsecase) CreatePost(v *validator.Validator, post *models.Post) (int, error) {
	post.Created = time.Now()
	errMap := validator.CreatePostValidation(post)

	postID, err := cpu.postRepository.CreatePost(post)
	if err != nil {
		return 0, err
	}

	for _, c := range post.Categories {
		category, err := cpu.categoryRepository.GetCategory(c)
		if err != nil {
			return 0, ErrCategoryNotFound
		}

		err = cpu.categoryRepository.AddCategory(post.ID, category.ID)
		if err != nil {
			return 0, err
		}
	}
	return postID, err
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
	v.Check(len(post.Title) <= 100, "title", "must not be more than 100 chars")

	v.Check(post.Content != "", "content", "must be provided")
	v.Check(len(post.Content) <= 100, "content", "must not be more than 100 chars")

}
