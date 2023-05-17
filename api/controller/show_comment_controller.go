package controller

import (
	"net/http"

	"forum/domain/usecase"
)

type CommentDetailsControler struct {
	CommentDetailsUsecase usecase.CommentDetailsUsecase
	*Controller
}

func (cdc *CommentDetailsControler) CommentDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		cdc.ResponseMethodNotAllowed(w)
	}

	commentID, err := GetIdFromURL2(2, r.URL.Path)
	if err != nil {
		cdc.ResponseNotFound(w)
	}

	comment, replies, err := cdc.CommentDetailsUsecase.GetComment(commentID)
	if err != nil {
		cdc.logger.PrintError(err)
		cdc.ResponseServerError(w)
		return
	}

	comment.Replies = replies

	cdc.Data.Comment = comment
	err = cdc.tmpl.ExecuteTemplate(w, "comment.html", cdc.Data)
	if err != nil {
		cdc.ResponseServerError(w)
		return
	}
}
