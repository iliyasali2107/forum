package usecase

import (
	"time"

	"forum/domain/models"
	"forum/domain/repository"
	"forum/pkg/validator"
)

type commentDetailsUsecase struct {
	commentRepository repository.CommentRepository
	ContextTimeout    time.Duration
}

type CommentDetailsUsecase interface {
	// GetCommentsByPostId(int) ([]*models.Comment, error)
	GetComment(int) (*models.Comment, []*models.Comment, error)
}

func NewCommentDetailsUsecase(commentRepositoy repository.CommentRepository, timeout time.Duration) CommentDetailsUsecase {
	return &commentDetailsUsecase{
		commentRepository: commentRepositoy,
		ContextTimeout:    timeout,
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

func (cdu *commentDetailsUsecase) GetComment(commentID int) (*models.Comment, []*models.Comment, error) {
	comment, err := cdu.commentRepository.GetComment(commentID)
	if err != nil {
		return nil, nil, err
	}

	replies, err := cdu.commentRepository.GetCommentReplies(commentID)
	if err != nil {
		return nil, nil, err
	}

	return comment, replies, nil
}

func validateComment(v *validator.Validator, comment *models.Comment) {
	v.Check(comment.Content != "", "content", "must be provided")
}
