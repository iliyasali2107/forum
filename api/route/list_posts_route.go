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

func NewListPostsRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	pr := repository.NewPostRepository(db)
	cr := repository.NewCategoryRepository(db)
	ur := repository.NewUserRepository(db)

	lpc := controller.ListPostsController{
		ListPostUsecase: usecase.NewListPostsUsecase(pr, cr, ur, timeout),
		Controller:      ctrl,
	}

	mux.HandleFunc(lpc.Data.Endpoints.PostsAllEndpoint, lpc.UserIdentity(lpc.ListPostsController))
}
