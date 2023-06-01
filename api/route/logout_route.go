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

func NewLogoutRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	ur := repository.NewUserRepository(db)
	lc := controller.LogoutConrtroller{
		LogoutUsecase: usecase.NewLogoutUsecase(ur, timeout),
		Controller:    ctrl,
	}

	mux.HandleFunc(lc.Data.Endpoints.LogoutEndpoint, lc.UserIdentity(lc.Authorized(lc.Logout)))
}
