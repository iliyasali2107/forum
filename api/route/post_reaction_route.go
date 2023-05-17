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

func NewPostReactionRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	pr := repository.NewPostRepository(db)
	ur := repository.NewUserRepository(db)
	cr := repository.NewCategoryRepository(db)
	rr := repository.NewReactionRepository(db)
	prc := controller.PostReactionController{
		PostReactionUsecase: usecase.NewPostReactionUsecase(rr, pr, cr, ur, timeout),
		Controller:          ctrl,
	}

	mux.HandleFunc(ctrl.Data.Endpoints.CreatePostReactionEndpoint, ctrl.UserIdentity(ctrl.Authorized(prc.PostReactionController)))
}
