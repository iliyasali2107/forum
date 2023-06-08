package route

import (
	"database/sql"
	"forum/api/controller"
	"forum/bootstrap"
	"forum/domain/repository"
	"forum/domain/usecase"
	"net/http"
	"time"
)

func NewCommentDetailsRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	cr := repository.NewCommentRepository(db)
	ur := repository.NewUserRepository(db)
	rr := repository.NewReactionRepository(db)
	cdc := controller.CommentDetailsControler{
		CommentDetailsUsecase: usecase.NewCommentDetailsUsecase(cr, ur, rr, timeout),
		Controller:            ctrl,
	}

	mux.HandleFunc(cdc.Data.Endpoints.CommentDetailsEndpoint, cdc.CommentDetails)
}
