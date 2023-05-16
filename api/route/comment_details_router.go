package route

import (
	"database/sql"
	"net/http"
	"time"

	"forum/api/controller"
	"forum/bootstrap"
	"forum/domain/repository"
	"forum/domain/usecase"
)

func NewCommentDetailsRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux) {
	cr := repository.NewCommentRepository(db)
	cdc := controller.CommentDetailsControler{
		CommentDetailsUsecase: usecase.NewCommentDetailsUsecase(cr, timeout),
	}

	mux.HandleFunc("/comment/", cdc.CommentDetails)
}
