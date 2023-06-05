package usecase

import (
	"fmt"
	"time"

	"forum/domain/models"
	"forum/domain/repository"
)

type createPostUsecase struct {
	postRepository     repository.PostRepository
	categoryRepository repository.CategoryRepository
	ContextTimeout     time.Duration
}

type CreatePostUsecase interface {
	CreatePost(*models.Post) (int, error)
	GetAllCategories() ([]*models.Category, error)
}

func NewCreatePostUsecase(postRepository repository.PostRepository, categoryRepository repository.CategoryRepository, timeout time.Duration) CreatePostUsecase {
	return &createPostUsecase{
		postRepository:     postRepository,
		categoryRepository: categoryRepository,
		ContextTimeout:     timeout,
	}
}

func (cpu *createPostUsecase) CreatePost(post *models.Post) (int, error) {
	post.Created = time.Now()

	postID, err := cpu.postRepository.CreatePost(post)
	if err != nil {
		return 0, fmt.Errorf("couldn't cerate post: %w", err)
	}

	for _, c := range post.Categories {
		category, err := cpu.categoryRepository.GetCategory(c)
		if err != nil {
			return 0, fmt.Errorf("couldn't get categories: %w", err)
		}

		err = cpu.categoryRepository.AddCategory(post.ID, category.ID)
		if err != nil {
			return 0, fmt.Errorf("couldn't add category: %w", err)
		}
	}
	return postID, nil
}

func (cpu *createPostUsecase) GetAllCategories() ([]*models.Category, error) {
	categories, err := cpu.categoryRepository.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("couldn't get all categories: %w", err)
	}

	return categories, nil
}
