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

func NewCommentCreateRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux) {
	cr := repository.NewCommentRepository(db)
	ur := repository.NewUserRepository(db)
	ccc := controller.CreateCommentController{
		CreateCommentUsecase: usecase.NewCreateCommentUsecase(cr, ur, timeout),
	}

	mux.HandleFunc("/posts/comment/create", ccc.CreateCommentController)
}
