package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/domain/models"
	"forum/domain/usecase"
)

type CreateCommentController struct {
	CreateCommentUsecase usecase.CreateCommentUsecase
	*Controller
}

func (pcc *CreateCommentController) CreateCommentController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		pcc.logger.PrintError(fmt.Errorf("method not allowed"))
		pcc.ResponseMethodNotAllowed(w)
		return
	}

	user := pcc.contextGetUser(r)
	if user != nil {
		pcc.Data.IsAuthorized = true
	} else {
		pcc.Data.IsAuthorized = false
	}
	
	err := r.ParseForm()
	if err != nil {
		pcc.logger.PrintError(fmt.Errorf("error while parsing request form: %w", err))
		pcc.ResponseBadRequest(w)
		return
	}

	comment := &models.Comment{}

	postIDStr := r.FormValue("post_id")
	postIDInt, err := strconv.Atoi(postIDStr)
	if err != nil {
		pcc.logger.PrintError(fmt.Errorf("couldn't convert post_id to int: %w", err))
		pcc.ResponseBadRequest(w)
		return
	}

	parentIDStr := r.FormValue("parent_id")
	parentIDInt, err := strconv.Atoi(parentIDStr)
	if err != nil {
		pcc.logger.PrintError(fmt.Errorf("couldn't convert parent_id to int: %w", err))
		pcc.ResponseBadRequest(w)
		return
	}

	comment.Content = strings.TrimSpace(r.FormValue("content"))
	if comment.Content == "" {
		pcc.logger.PrintError(fmt.Errorf("empty comment content is not acceptable"))
		pcc.ResponseBadRequest(w)
		return
	}
	comment.UserID = user.ID
	comment.PostID = postIDInt
	comment.ParentID = parentIDInt

	err = pcc.CreateCommentUsecase.CreateComment(comment)
	if err != nil {
		pcc.logger.PrintError(err)
		pcc.ResponseServerError(w)
		return
	}

	http.Redirect(w, r, pcc.Data.Endpoints.PostDetailsEndpoint+strconv.Itoa(postIDInt), http.StatusSeeOther)
}
