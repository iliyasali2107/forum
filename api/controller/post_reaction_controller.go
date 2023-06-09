package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/domain/models"
	"forum/domain/usecase"
)

type PostReactionController struct {
	PostReactionUsecase usecase.PostReactionUsecase
	*Controller
}

func (prc *PostReactionController) PostReactionController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		prc.logger.PrintError(fmt.Errorf("Controller: LikePost: method not allowed"))
		prc.ResponseMethodNotAllowed(w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		prc.logger.PrintError(fmt.Errorf("error while parsing form: %w", err))
		prc.ResponseBadRequest(w)
		return
	}

	postIDStr := r.FormValue("post_id")
	postIDInt, err := strconv.Atoi(postIDStr)
	if err != nil {
		prc.logger.PrintError(fmt.Errorf("error while converting post_id to int: %w", err))
		prc.ResponseBadRequest(w)
		return
	}

	reactionTypeStr := r.FormValue("reaction")
	reactionTypeInt, err := strconv.Atoi(reactionTypeStr)
	if err != nil || (reactionTypeInt != 1 && reactionTypeInt != 0) {
		prc.logger.PrintError(fmt.Errorf("incorrect reaction type: %w", err))
		prc.ResponseBadRequest(w)
		return
	}

	post, err := prc.PostReactionUsecase.GetPost(postIDInt)
	if err != nil {
		prc.logger.PrintError(fmt.Errorf("Controller: LikePost: GetPost: %w", err))
		prc.ResponseServerError(w)
		return
	}

	user := prc.contextGetUser(r)
	if user != nil {
		prc.Data.IsAuthorized = true
	} else {
		prc.Data.IsAuthorized = false
	}

	reaction := &models.Reaction{UserID: user.ID, PostID: postIDInt, Type: reactionTypeInt}
	if reaction.Type == 1 {
		err = prc.PostReactionUsecase.LikePost(reaction)
	} else if reaction.Type == 0 {
		err = prc.PostReactionUsecase.DislikePost(reaction)
	}

	if err != nil {
		prc.logger.PrintError(err)
		prc.ResponseServerError(w)
		return
	}

	http.Redirect(w, r, prc.Data.Endpoints.PostDetailsEndpoint+strconv.Itoa(post.ID), http.StatusSeeOther)
}
