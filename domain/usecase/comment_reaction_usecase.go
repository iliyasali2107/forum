package usecase

import (
	"database/sql"
	"fmt"
	"time"

	"forum/domain/models"
	"forum/domain/repository"
)

type commentReactionUsecase struct {
	reactionRepostitory repository.ReactionRepository
	ContextTimeout      time.Duration
}

type CommentReactionUsecase interface {
	LikeComment(reaction *models.Reaction) error
	DislikeComment(reaction *models.Reaction) error
	CommentLikeCount(commentID int) (int, error)
	CommentDislikeCount(commentID int) (int, error)
}

func NewCommentReactionUsecase(reactionRepository repository.ReactionRepository, timeout time.Duration) CommentReactionUsecase {
	return &commentReactionUsecase{
		reactionRepostitory: reactionRepository,
		ContextTimeout:      timeout,
	}
}

func (cru *commentReactionUsecase) LikeComment(reaction *models.Reaction) error {
	dbReaction, err := cru.reactionRepostitory.GetCommentReaction(reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := cru.reactionRepostitory.CreateCommentReaction(reaction)
			if err != nil {
				return fmt.Errorf("couldn't create reaction to comment: %w", err)
			}

			return nil
		}

		return fmt.Errorf("error while trying to get dbReaction: %w", err)
	}

	switch dbReaction.Type {
	case 1:
		if err := cru.reactionRepostitory.DeleteCommentReaction(dbReaction); err != nil {
			return fmt.Errorf("couldn't delete reaction: %w", err)
		}
	case 0:
		if err := cru.reactionRepostitory.UpdateCommentReaction(reaction); err != nil {
			return fmt.Errorf("couldn't update reaction: %w", err)
		}
	}

	return nil
}

func (cru *commentReactionUsecase) DislikeComment(reaction *models.Reaction) error {
	dbReaction, err := cru.reactionRepostitory.GetCommentReaction(reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := cru.reactionRepostitory.CreateCommentReaction(reaction)
			if err != nil {
				return fmt.Errorf("couldn't create reaction to comment: %w", err)
			}

			return nil
		}

		return fmt.Errorf("error while trying to get dbReaction: %w", err)

	}
	switch dbReaction.Type {
	case 1:
		if err := cru.reactionRepostitory.UpdateCommentReaction(reaction); err != nil {
			return fmt.Errorf("couldn't update reaction: %w", err)
		}
	case 0:
		if err := cru.reactionRepostitory.DeleteCommentReaction(dbReaction); err != nil {
			return fmt.Errorf("couldn't delete reaction: %w", err)
		}
	}

	return nil
}

func (cru *commentReactionUsecase) CommentLikeCount(commentID int) (int, error) {
	likes, err := cru.reactionRepostitory.GetCommentLikes(commentID)
	if err != nil {
		return 0, fmt.Errorf("couldn't get comment likes: %w", err)
	}

	return likes, nil
}

func (cru *commentReactionUsecase) CommentDislikeCount(commentID int) (int, error) {
	dislikes, err := cru.reactionRepostitory.GetCommentDislikes(commentID)
	if err != nil {
		return 0, fmt.Errorf("couldn't get comment dislikes: %w", err)
	}

	return dislikes, nil
}
