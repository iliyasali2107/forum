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

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	ur := repository.NewUserRepository(db)
	lc := controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(ur, env, timeout),
		Controller:   ctrl,
	}

	mux.HandleFunc(lc.Data.Endpoints.LoginEndpoint, lc.UserIdentity(lc.NotAuthorized(lc.Login)))
}
