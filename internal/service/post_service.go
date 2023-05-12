package service

import (
	"time"

	"forum/internal/models"
	"forum/internal/repository"
	"forum/pkg/validator"
)

type PostService interface {
	CreatePost(*validator.Validator, *models.Post) error
	GetPost(id int) (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
	GetLikedPosts(userID int) ([]*models.Post, error)
	GetDislikedPosts(userID int) ([]*models.Post, error)
	GetCreatedPosts(userID int) ([]*models.Post, error)
	GetCommentedPosts(userID int) ([]*models.Post, error)
	GetAllCategories() ([]*models.Category, error)
}

type postService struct {
	pr repository.PostRepository
	cr repository.CategoryRepository
	ur repository.UserRepository
}

func NewPostService(postRepository repository.PostRepository, categoryRepository repository.CategoryRepository, userRepository repository.UserRepository) PostService {
	return &postService{
		pr: postRepository,
		cr: categoryRepository,
		ur: userRepository,
	}
}

func (ps *postService) CreatePost(v *validator.Validator, post *models.Post) error {
	post.Created = time.Now()
	if validatePost(v, post); !v.Valid() {
		return ErrFormValidation
	}

	_, err := ps.pr.CreatePost(post)
	if err != nil {
		return err
	}

	for _, c := range post.Categories {
		category, err := ps.cr.GetCategory(c)
		if err != nil {
			return ErrCategoryNotFound
		}

		err = ps.cr.AddCategory(post.ID, category.ID)
		if err != nil {
			return err
		}
	}
	return err
}

func (ps *postService) GetAllCategories() ([]*models.Category, error) {
	categories, err := ps.cr.GetAllCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (ps *postService) GetPost(id int) (*models.Post, error) {
	post, err := ps.pr.GetPost(id)
	if err != nil {
		return nil, err
	}

	user, err := ps.ur.GetUser(post.User.ID)
	if err != nil {
		return nil, err
	}

	post.User = user

	err = ps.cr.GetCategoriesForPost(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *postService) GetAllPosts() ([]*models.Post, error) {
	posts, err := ps.pr.GetAllPosts()
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		user, err := ps.ur.GetUser(post.User.ID)
		if err != nil {
			return nil, err
		}
		post.User = user
		err = ps.cr.GetCategoriesForPost(post)
		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

func (ps *postService) GetCreatedPosts(userID int) ([]*models.Post, error) {
	posts, err := ps.pr.GetAllPostsOfUser(userID)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		user, err := ps.ur.GetUser(post.User.ID)
		if err != nil {
			return nil, err
		}
		post.User = user
		err = ps.cr.GetCategoriesForPost(post)
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func (ps *postService) GetLikedPosts(userID int) ([]*models.Post, error) {
	return nil, nil
}

func (ps *postService) GetDislikedPosts(userID int) ([]*models.Post, error) {
	return nil, nil
}

func (ps *postService) GetCommentedPosts(userID int) ([]*models.Post, error) {
	return nil, nil
}

func validatePost(v *validator.Validator, post *models.Post) {
	v.Check(post.Title != "", "title", "must be provided")
	v.Check(len(post.Title) <= 20, "title", "must not be more than 20 chars")

	v.Check(post.Content != "", "content", "must be provided")
	v.Check(len(post.Content) <= 20, "content", "must not be more than 20 chars")
}
