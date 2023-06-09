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

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	ur := repository.NewUserRepository(db)
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, timeout),
		Controller:    ctrl,
	}

	mux.HandleFunc(sc.Data.Endpoints.SignupEndpoint, sc.UserIdentity(sc.NotAuthorized(sc.Signup)))
}
