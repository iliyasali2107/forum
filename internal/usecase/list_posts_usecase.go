package usecase

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type listPostsUsecase struct {
	postRepository     repository.PostRepository
	categoryRepository repository.CategoryRepository
	userRepository     repository.UserRepository
}

type ListPostsUsecase interface {
	GetAllPosts() ([]*models.Post, error)
	GetLikedPosts(userID int) ([]*models.Post, error)
	GetDislikedPosts(userID int) ([]*models.Post, error)
	GetCreatedPosts(userID int) ([]*models.Post, error)
	GetCommentedPosts(userID int) ([]*models.Post, error)
}

func NewListPostsUsecase(postRepository repository.PostRepository, categoryRepository repository.CategoryRepository, userRepository repository.UserRepository) ListPostsUsecase {
	return &listPostsUsecase{
		postRepository:     postRepository,
		categoryRepository: categoryRepository,
		userRepository:     userRepository,
	}
}

func (lpu *listPostsUsecase) GetAllPosts() ([]*models.Post, error) {
	posts, err := lpu.postRepository.GetAllPosts()
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		user, err := lpu.userRepository.GetUser(post.User.ID)
		if err != nil {
			return nil, err
		}
		post.User = user
		err = lpu.categoryRepository.GetCategoriesForPost(post)
		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func (lpu *listPostsUsecase) GetLikedPosts(userID int) ([]*models.Post, error) {
	return nil, nil
}

func (lpu *listPostsUsecase) GetDislikedPosts(userID int) ([]*models.Post, error) {
	return nil, nil
}

func (lpu *listPostsUsecase) GetCreatedPosts(userID int) ([]*models.Post, error) {
	return nil, nil
}

func (lpu *listPostsUsecase) GetCommentedPosts(userID int) ([]*models.Post, error) {
	return nil, nil
}
