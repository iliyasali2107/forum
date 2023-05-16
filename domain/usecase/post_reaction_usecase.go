package usecase

import (
	"database/sql"
	"time"

	"forum/domain/models"
	"forum/domain/repository"
)

type postReactionUsecase struct {
	reactionRepository repository.ReactionRepository
	postRepository     repository.PostRepository
	userRepository     repository.UserRepository
	categoryRepository repository.CategoryRepository
	ContextTimeout     time.Duration
}

type PostReactionUsecase interface {
	GetPost(postID int) (*models.Post, error)
	LikePost(reaction *models.Reaction) error
	DislikePost(reaction *models.Reaction) error
	GetPostLikes(post_id int) (int, error)
	GetPostDislikes(post_id int) (int, error)
}

func NewPostReactionUsecase(reactionRepository repository.ReactionRepository, postRepository repository.PostRepository, categoryRepository repository.CategoryRepository, userRepository repository.UserRepository, timeout time.Duration) PostReactionUsecase {
	return &postReactionUsecase{
		reactionRepository: reactionRepository,
		postRepository:     postRepository,
		userRepository:     userRepository,
		categoryRepository: categoryRepository,
		ContextTimeout:     timeout,
	}
}

func (pru *postReactionUsecase) GetPost(postID int) (*models.Post, error) {
	post, err := pru.postRepository.GetPost(postID)
	if err != nil {
		return nil, err
	}

	user, err := pru.userRepository.GetUser(post.User.ID)
	if err != nil {
		return nil, err
	}

	post.User = user

	err = pru.categoryRepository.GetCategoriesForPost(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pru *postReactionUsecase) LikePost(reaction *models.Reaction) error {
	dbReaction, err := pru.reactionRepository.GetPostReaction(reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			err := pru.reactionRepository.CreatePostReaction(reaction)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	switch dbReaction.Type {
	case 1:
		if err := pru.reactionRepository.DeletePostReaction(dbReaction); err != nil {
			return err
		}
	case 0:
		if err := pru.reactionRepository.UpdatePostReaction(reaction); err != nil {
			return err
		}
	}
	return nil
}

func (pru *postReactionUsecase) DislikePost(reaction *models.Reaction) error {
	dbReaction, err := pru.reactionRepository.GetPostReaction(reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			err = pru.reactionRepository.CreatePostReaction(reaction)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	switch dbReaction.Type {
	case 1:
		if err := pru.reactionRepository.UpdatePostReaction(reaction); err != nil {
			return err
		}
	case 0:
		if err := pru.reactionRepository.DeletePostReaction(dbReaction); err != nil {
			return err
		}
	}
	return nil
}

func (pru *postReactionUsecase) GetPostLikes(post_id int) (int, error) {
	likes, err := pru.reactionRepository.GetPostLikes(post_id)
	if err != nil {
		return 0, err
	}

	return likes, nil
}

func (pru *postReactionUsecase) GetPostDislikes(post_id int) (int, error) {
	dislikes, err := pru.reactionRepository.GetPostDislikes(post_id)
	if err != nil {
		return 0, err
	}

	return dislikes, nil
}
