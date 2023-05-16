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

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	ur := repository.NewUserRepository(db)
	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, timeout),
		Controller:    ctrl,
	}

	mux.HandleFunc("/signup", sc.Signup)
}
