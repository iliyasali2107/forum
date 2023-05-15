package controller

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/service"
	"forum/internal/usecase"
	"net/http"
	"strconv"
)

type PostCommentController struct {
	PostCommentUsecase usecase.PostCommentUsecase
	Controller
}

func (pcc *PostCommentController) CreateCommentController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		user := pcc.contextGetUser(r)

		err := r.ParseForm()
		if err != nil {
			pcc.logger.PrintError(fmt.Errorf("Controller: comment-create: ParseForm"))
			pcc.logger.PrintError(err)
			pcc.ResponseBadRequest(w)
			return
		}

		comment := &models.Comment{}

		postIDStr := r.FormValue("post_id")
		postIDInt, err := strconv.Atoi(postIDStr)
		if err != nil {
			pcc.logger.PrintError(fmt.Errorf("Controller: comment-create: strconv post_id"))
			pcc.logger.PrintError(err)
			pcc.ResponseBadRequest(w)
		}

		parentIDStr := r.FormValue("parent_id")
		parentIDInt, err := strconv.Atoi(parentIDStr)
		if err != nil {
			pcc.logger.PrintError(fmt.Errorf("Controller: comment-create: strconv parent_id"))
			pcc.logger.PrintError(err)
			pcc.ResponseBadRequest(w)
			return
		}

		comment.Content = r.FormValue("content")
		comment.UserID = user.ID
		comment.PostID = postIDInt
		comment.ParentID = parentIDInt

		err = pcc.PostCommentUsecase.CreateComment(pcc.validator, comment)
		if err != nil {
			if err == service.ErrFormValidation {
				pcc.logger.PrintError(fmt.Errorf("Controller: comment-create: CreateComment ErrFormValidation"))
				pcc.logger.PrintError(err)
				pcc.ResponseBadRequest(w)
				return
			}
			pcc.logger.PrintError(fmt.Errorf("Controller: comment-create: CreateComment"))
			pcc.logger.PrintError(err)
			pcc.ResponseServerError(w)
			return
		}

		http.Redirect(w, r, "/posts/"+strconv.Itoa(postIDInt), http.StatusSeeOther)
	default:
		pcc.logger.PrintError(fmt.Errorf("Controller: comment-create: method not allowed"))
		pcc.ResponseMethodNotAllowed(w)
		return
	}
}
