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

func NewCommentReactionRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	rr := repository.NewReactionRepository(db)

	ccrr := controller.CommentReactionController{
		CommentReactionUsecase: usecase.NewCommentReactionUsecase(rr, timeout),
		Controller:             ctrl,
	}

	mux.HandleFunc(ccrr.Data.Endpoints.CommentReactionEndpoint, ccrr.UserIdentity(ccrr.Authorized(ccrr.CommentReactionController)))
}
