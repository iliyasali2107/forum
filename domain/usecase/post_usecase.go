package usecase

import (
	"fmt"
	"forum/domain/models"
	"forum/domain/repository"
	"time"
)

type postDetailsUsecase struct {
	postRepository     repository.PostRepository
	userRepository     repository.UserRepository
	categoryRepository repository.CategoryRepository
	reactionRepository repository.ReactionRepository
	commentRepository  repository.CommentRepository
	ContextTimeout     time.Duration
}

type PostDetailsUsecase interface {
	GetPost(postID int) (*models.Post, error)
	GetPostLikes(postID int) (int, error)
	GetPostDislikes(postID int) (int, error)
	GetCommentsByPostId(postID int) ([]*models.Comment, error)
}

func NewPostDetailsUsecase(postRepository repository.PostRepository, userRepository repository.UserRepository, categoryRepository repository.CategoryRepository, reactionRepository repository.ReactionRepository, commentRepository repository.CommentRepository, timeout time.Duration) PostDetailsUsecase {
	return &postDetailsUsecase{
		postRepository:     postRepository,
		userRepository:     userRepository,
		categoryRepository: categoryRepository,
		reactionRepository: reactionRepository,
		commentRepository:  commentRepository,
		ContextTimeout:     timeout,
	}
}

func (pdu *postDetailsUsecase) GetPost(postID int) (*models.Post, error) {
	post, err := pdu.postRepository.GetPost(postID)
	if err != nil {
		return nil, err
	}

	post.CreatedStr = post.Created.Format("2006-01-02 15:04:05")

	user, err := pdu.userRepository.GetUser(post.User.ID)
	if err != nil {
		return nil, err
	}

	post.User = user

	err = pdu.categoryRepository.GetCategoriesForPost(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pdu *postDetailsUsecase) GetPostLikes(postID int) (int, error) {
	likes, err := pdu.reactionRepository.GetPostLikes(postID)
	if err != nil {
		return 0, err
	}

	return likes, nil
}

func (pdu *postDetailsUsecase) GetPostDislikes(postID int) (int, error) {
	dislikes, err := pdu.reactionRepository.GetPostDislikes(postID)
	if err != nil {
		return 0, err
	}

	return dislikes, nil
}

func (pdu *postDetailsUsecase) GetCommentsByPostId(post_id int) ([]*models.Comment, error) {
	comments, err := pdu.commentRepository.GetPostComments(post_id)
	if err != nil {
		return nil, fmt.Errorf("couldn't get post comments: %w", err)
	}

	for _, comment := range comments {
		replies, err := pdu.commentRepository.GetCommentReplies(comment.ID)
		if err != nil {
			return nil, fmt.Errorf("couldn't get comment repiles count: %w", err)
		}

		for _, reply := range replies {
			replyUser, err := pdu.userRepository.GetUser(reply.UserID)
			if err != nil {
				return nil, fmt.Errorf("couldn't get user: %w", err)
			}
			reply.User = replyUser

			likes, err := pdu.reactionRepository.GetCommentLikes(reply.ID)
			if err != nil {
				return nil, fmt.Errorf("couldn't get reply likes: %w", err)
			}

			reply.Likes = likes

			dislikes, err := pdu.reactionRepository.GetCommentDislikes(reply.ID)
			if err != nil {
				return nil, fmt.Errorf("couldn't get reply dislikes: %w", err)
			}

			reply.Dislikes = dislikes

		}

		comment.Replies = replies

		user, err := pdu.userRepository.GetUser(comment.UserID)
		if err != nil {
			return nil, fmt.Errorf("couldn't get user: %w", err)
		}

		comment.User = user
	}

	return comments, nil
}
