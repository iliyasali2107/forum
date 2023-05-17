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

func NewCommentCreateRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	cr := repository.NewCommentRepository(db)
	ur := repository.NewUserRepository(db)
	ccc := controller.CreateCommentController{
		CreateCommentUsecase: usecase.NewCreateCommentUsecase(cr, ur, timeout),
		Controller:           ctrl,
	}

	mux.HandleFunc(ctrl.Data.Endpoints.CreateCommentEndpoint, ctrl.UserIdentity(ctrl.Authorized(ccc.CreateCommentController)))
}
