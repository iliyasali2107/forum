package usecase

import (
	"database/sql"
	"fmt"
	"time"

	"forum/domain/models"
	"forum/domain/repository"
)

type commentDetailsUsecase struct {
	commentRepository  repository.CommentRepository
	userRepository     repository.UserRepository
	reactionRepository repository.ReactionRepository
	ContextTimeout     time.Duration
}

type CommentDetailsUsecase interface {
	// GetCommentsByPostId(int) ([]*models.Comment, error)
	GetComment(int) (*models.Comment, []*models.Comment, error)
	CommentLikeCount(commentID int) (int, error)
	CommentDislikeCount(commentID int) (int, error)
}

func NewCommentDetailsUsecase(commentRepositoy repository.CommentRepository, userRepository repository.UserRepository, reactionRepository repository.ReactionRepository, timeout time.Duration) CommentDetailsUsecase {
	return &commentDetailsUsecase{
		commentRepository:  commentRepositoy,
		userRepository:     userRepository,
		reactionRepository: reactionRepository,
		ContextTimeout:     timeout,
	}
}

// func (cdu *commentDetailsUsecase) GetCommentsByPostId(post_id int) ([]*models.Comment, error) {
// 	comments, err := cdu.commentRepository.GetPostComments(post_id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, comment := range comments {
// 		replies, err := cdu.commentRepository.GetCommentRepliesCount(comment.ID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		comment.ReplyCount = replies
// 	}

// 	return comments, nil
// }

func (cdu *commentDetailsUsecase) CommentLikeCount(commentID int) (int, error) {
	likes, err := cdu.reactionRepository.GetCommentLikes(commentID)
	if err != nil {
		return 0, fmt.Errorf("couldn't get comment likes: %w", err)
	}

	return likes, nil
}

func (cdu *commentDetailsUsecase) CommentDislikeCount(commentID int) (int, error) {
	dislikes, err := cdu.reactionRepository.GetCommentDislikes(commentID)
	if err != nil {
		return 0, fmt.Errorf("couldn't get comment dislikes: %w", err)
	}

	return dislikes, nil
}

func (cdu *commentDetailsUsecase) GetComment(commentID int) (*models.Comment, []*models.Comment, error) {
	comment, err := cdu.commentRepository.GetComment(commentID)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't get comment: %w", err)
	}

	likes, err := cdu.reactionRepository.GetCommentLikes(comment.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, nil, fmt.Errorf("couldn't get comment likes: %w", err)
		}
	}

	comment.Likes = likes

	dislikes, err := cdu.reactionRepository.GetCommentDislikes(comment.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, nil, fmt.Errorf("couldn't get comment dislikes: %w", err)
		}
	}

	comment.Dislikes = dislikes

	replies, err := cdu.commentRepository.GetCommentReplies(commentID)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't get comment replies: %w", err)
	}

	for _, reply := range replies {
		user, err := cdu.userRepository.GetUser(reply.UserID)
		if err != nil {
			return nil, nil, fmt.Errorf("couldn't get user: %w", err)
		}

		reply.User = user

		likes, err := cdu.reactionRepository.GetCommentLikes(reply.ID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, nil, fmt.Errorf("couldn't get reply likes: %w", err)
			}
		}

		reply.Likes = likes

		dislikes, err := cdu.reactionRepository.GetCommentDislikes(reply.ID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, nil, fmt.Errorf("couldn't get reply dislikes: %w", err)
			}
		}

		reply.Dislikes = dislikes

	}

	return comment, replies, nil
}
