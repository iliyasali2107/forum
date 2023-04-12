package service

import (
	"database/sql"
	"forum/internal/models"
	"forum/internal/repository"
)

type ReactionService interface {
	LikePost(reaction *models.Reaction) error
	DislikePost(reaction *models.Reaction) error
}

type reactionService struct {
	rr repository.ReactionRepository
}

func NewReactionService(repository repository.ReactionRepository) ReactionService {
	return &reactionService{
		rr: repository,
	}
}

func (rs *reactionService) LikePost(reaction *models.Reaction) error {
	dbReaction, err := rs.rr.GetReactionOfPost(reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			err := rs.rr.CreatePostReaction(reaction)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	switch dbReaction.Type {
	case 1:
		if err := rs.rr.DeletePostReaction(dbReaction); err != nil {
			return err
		}
	case 0:
		if err := rs.rr.DeletePostReaction(dbReaction); err != nil {
			return err
		}

		if err := rs.rr.CreatePostReaction(reaction); err != nil {
			return err
		}
	}
	return nil
}

func (rs *reactionService) DislikePost(reaction *models.Reaction) error {
	dbReaction, err := rs.rr.GetReactionOfPost(reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			err := rs.rr.CreatePostReaction(reaction)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	switch dbReaction.Type {
	case 1:
		if err := rs.rr.DeletePostReaction(dbReaction); err != nil {
			return err
		}

		if err := rs.rr.CreatePostReaction(reaction); err != nil {
			return err
		}

	case 0:
		if err := rs.rr.DeletePostReaction(dbReaction); err != nil {
			return err
		}
	}
	return nil
}
