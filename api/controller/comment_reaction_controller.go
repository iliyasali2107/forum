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
		crc.logger.PrintError(fmt.Errorf("Controller: LikeComment: method not allowed"))
		crc.ResponseMethodNotAllowed(w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		crc.logger.PrintError(err)
		crc.ResponseBadRequest(w)
		return
	}

	commentIDStr := r.FormValue("comment_id")
	commentId, err := strconv.Atoi(commentIDStr)
	if err != nil {
		crc.logger.PrintError(err)
		crc.ResponseBadRequest(w)
		return
	}

	reactionTypeStr := r.FormValue("reaction")
	reactionType, err := strconv.Atoi(reactionTypeStr)
	if err != nil || (reactionType != 1 && reactionType != 0) {
		crc.logger.PrintError(err)
		crc.ResponseBadRequest(w)
		return
	}

	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)

	user := crc.contextGetUser(r)

	reaction := &models.Reaction{
		UserID:    user.ID,
		CommentID: commentId,
		Type:      reactionType,
	}

	if reaction.Type == 1 {
		err = crc.CommentReactionUsecase.LikeComment(reaction)
	} else {
		err = crc.CommentReactionUsecase.DislikeComment(reaction)
	}

	if err != nil {
		crc.logger.PrintError(err)
		crc.ResponseServerError(w)
		return
	}

	http.Redirect(w, r, crc.Data.Endpoints.PostDetailsEndpoint+strconv.Itoa(postID), http.StatusSeeOther)
}
