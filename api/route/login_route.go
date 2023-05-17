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

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	ur := repository.NewUserRepository(db)
	lc := controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(ur, env, timeout),
		Controller:   ctrl,
	}

	mux.HandleFunc(ctrl.Data.Endpoints.LoginEndpoint, lc.Login)
}
