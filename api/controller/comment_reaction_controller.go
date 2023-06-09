package controller

import (
	"fmt"
	"forum/domain/models"
	"forum/domain/usecase"
	"net/http"
	"strconv"
)

type CommentReactionController struct {
	CommentReactionUsecase usecase.CommentReactionUsecase
	*Controller
}

func (crc *CommentReactionController) CommentReactionController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		crc.logger.PrintError(fmt.Errorf("comment_reaction: method not allowed"))
		crc.ResponseMethodNotAllowed(w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		crc.logger.PrintError(fmt.Errorf("error while parsing request form: %w", err))
		crc.ResponseBadRequest(w)
		return
	}

	commentIDStr := r.FormValue("comment_id")
	commentId, err := strconv.Atoi(commentIDStr)
	if err != nil {
		crc.logger.PrintError(fmt.Errorf("incorrect comment_id value: %w", err))
		crc.ResponseBadRequest(w)
		return
	}

	reactionTypeStr := r.FormValue("reaction")
	reactionType, err := strconv.Atoi(reactionTypeStr)
	if err != nil || (reactionType != 1 && reactionType != 0) {
		crc.logger.PrintError(fmt.Errorf("incorrect reactionType value"))
		crc.ResponseBadRequest(w)
		return
	}

	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		crc.logger.PrintError(fmt.Errorf("incorrect post_id value: %w", err))
		crc.ResponseBadRequest(w)
		return
	}

	user := crc.contextGetUser(r)
	if user != nil {
		crc.Data.IsAuthorized = true
	} else {
		crc.Data.IsAuthorized = false
	}

	reaction := &models.Reaction{
		UserID:    user.ID,
		CommentID: commentId,
		Type:      reactionType,
	}
	if reaction.Type == 1 {
		err = crc.CommentReactionUsecase.LikeComment(reaction)
	} else if reaction.Type == 0 {
		err = crc.CommentReactionUsecase.DislikeComment(reaction)
	}

	if err != nil {
		crc.logger.PrintError(fmt.Errorf("comment-reaction: %w", err))
		crc.ResponseServerError(w)
		return
	}

	nextPage := r.FormValue("next")
	if nextPage == "" {
		http.Redirect(w, r, crc.Data.Endpoints.PostDetailsEndpoint+strconv.Itoa(postID), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, nextPage, http.StatusSeeOther)
	}
}
