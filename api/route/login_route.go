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

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux) {
	ur := repository.NewUserRepository(db)
	lc := controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(ur, env, timeout),
	}

	mux.HandleFunc("/login", lc.Login)
}
