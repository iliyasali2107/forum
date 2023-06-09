package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/domain/models"
	"forum/domain/repository"
	"forum/pkg/utils"
	"time"
)

type listPostsUsecase struct {
	postRepository     repository.PostRepository
	categoryRepository repository.CategoryRepository
	userRepository     repository.UserRepository
	ContextTimeout     time.Duration
}

type ListPostsUsecase interface {
	GetAllPosts() ([]*models.Post, error)
	GetLikedPosts(userID int) ([]*models.Post, error)
	GetDislikedPosts(userID int) ([]*models.Post, error)
	GetCreatedPosts(userID int) ([]*models.Post, error)
	GetCommentedPosts(userID int) ([]*models.Post, error)
	GetAllCategories() ([]*models.Category, error)
	GetPostsByCategories(ids ...int) ([]*models.Post, error)
}

func NewListPostsUsecase(postRepository repository.PostRepository, categoryRepository repository.CategoryRepository, userRepository repository.UserRepository, timeout time.Duration) ListPostsUsecase {
	return &listPostsUsecase{
		postRepository:     postRepository,
		categoryRepository: categoryRepository,
		userRepository:     userRepository,
		ContextTimeout:     timeout,
	}
}

func (lpu *listPostsUsecase) GetAllPosts() ([]*models.Post, error) {
	posts, err := lpu.postRepository.GetAllPosts()
	if err != nil {
		return nil, fmt.Errorf("couldn't get all posts: %w", err)
	}

	for _, post := range posts {
		user, err := lpu.userRepository.GetUser(post.User.ID)
		if err != nil {
			return nil, fmt.Errorf("couldn't get user: %w", err)
		}
		post.User = user
		post.CreatedStr = post.Created.Format("2006-01-02 15:04:05")
		err = lpu.categoryRepository.GetCategoriesForPost(post)
		if err != nil {
			return nil, fmt.Errorf("couldn't get categories of post: %w", err)
		}
	}

	return posts, nil
}

func (lpu *listPostsUsecase) GetLikedPosts(userID int) ([]*models.Post, error) {
	posts, err := lpu.postRepository.GetLikedPosts(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNoPosts
		}
		return nil, fmt.Errorf("couldn't get user's liked posts: %w", err)
	}

	for _, post := range posts {
		user, err := lpu.userRepository.GetUser(post.User.ID)
		if err != nil {
			return nil, fmt.Errorf("couldn't get user: %w", err)
		}

		post.User = user
		err = lpu.categoryRepository.GetCategoriesForPost(post)
		if err != nil {
			return nil, fmt.Errorf("couldn't get categories of post: %w", err)
		}
		post.CreatedStr = post.Created.Format("2006-01-02 15:04:05")
	}

	return posts, nil
}

func (lpu *listPostsUsecase) GetDislikedPosts(userID int) ([]*models.Post, error) {
	posts, err := lpu.postRepository.GetDislikedPosts(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNoPosts
		}
		return nil, fmt.Errorf("couldn't get user's disliked posts: %w", err)
	}

	for _, post := range posts {
		err = lpu.categoryRepository.GetCategoriesForPost(post)
		if err != nil {
			return nil, fmt.Errorf("couldn't get categories of post: %w", err)
		}
		post.CreatedStr = post.Created.Format("2006-01-02 15:04:05")
	}

	return posts, nil
}

func (lpu *listPostsUsecase) GetCreatedPosts(userID int) ([]*models.Post, error) {
	posts, err := lpu.postRepository.GetAllPostsOfUser(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNoPosts
		}
		return nil, fmt.Errorf("couldn't get users posts: %w", err)
	}

	for _, post := range posts {
		user, err := lpu.userRepository.GetUser(post.User.ID)
		if err != nil {
			return nil, fmt.Errorf("couldn't get user: %w", err)
		}
		post.User = user
		err = lpu.categoryRepository.GetCategoriesForPost(post)
		if err != nil {
			return nil, fmt.Errorf("couldn't get categories of post: %w", err)
		}
		post.CreatedStr = post.Created.Format("2006-01-02 15:04:05")
	}

	return posts, nil
}

func (lpu *listPostsUsecase) GetPostsByCategories(ids ...int) ([]*models.Post, error) {
	posts, err := lpu.postRepository.GetPostsByCategory(ids...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNoPosts
		}
		return nil, fmt.Errorf("couldn't get posts: %w", err)
	}

	for _, post := range posts {
		user, err := lpu.userRepository.GetUser(post.User.ID)
		if err != nil {
			return nil, fmt.Errorf("couldn't get user: %w", err)
		}
		post.User = user
		err = lpu.categoryRepository.GetCategoriesForPost(post)
		if err != nil {
			return nil, fmt.Errorf("couldn't get categories of post: %w", err)
		}
		post.CreatedStr = post.Created.Format("2006-01-02 15:04:05")
	}

	return posts, nil
}

func (lpu *listPostsUsecase) GetCommentedPosts(userID int) ([]*models.Post, error) {
	return nil, nil
}

func (lpu *listPostsUsecase) GetAllCategories() ([]*models.Category, error) {
	categories, err := lpu.categoryRepository.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("couldn't get categories: %w", err)
	}

	return categories, nil
}
