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
		prc.logger.PrintError(err)
		prc.ResponseBadRequest(w)
		return
	}

	postIDStr := r.FormValue("post_id")
	postIDInt, err := strconv.Atoi(postIDStr)
	if err != nil {
		prc.logger.PrintError(err)
		prc.ResponseBadRequest(w)
		return
	}

	reactionTypeStr := r.FormValue("reaction")
	reactionTypeInt, err := strconv.Atoi(reactionTypeStr)
	if err != nil || (reactionTypeInt != 1 && reactionTypeInt != 0) {
		prc.logger.PrintError(err)
		prc.ResponseBadRequest(w)
		return
	}

	post, err := prc.PostReactionUsecase.GetPost(postIDInt)
	if err != nil {
		prc.logger.PrintError(fmt.Errorf("Controller: LikePost: GetPost"))
		prc.logger.PrintError(err)
		prc.ResponseServerError(w)
		return
	}
	user := prc.contextGetUser(r)

	reaction := &models.Reaction{UserID: user.ID, PostID: postIDInt, Type: reactionTypeInt}

	if reactionTypeInt == 1 {
		err = prc.PostReactionUsecase.LikePost(reaction)
	} else if reactionTypeInt == 0 {
		err = prc.PostReactionUsecase.DislikePost(reaction)
	}

	if err != nil {
		prc.logger.PrintError(fmt.Errorf("Controller: LikePost: prc.PostReactionUsecase.LikePost()"))
		prc.logger.PrintError(err)
		prc.ResponseServerError(w)
		return
	}

	http.Redirect(w, r, prc.Data.Endpoints.PostDetailsEndpoint+strconv.Itoa(post.ID), http.StatusSeeOther)
}
