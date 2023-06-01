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

func NewCreatePostRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	pr := repository.NewPostRepository(db)
	cr := repository.NewCategoryRepository(db)
	cpc := controller.CreatePostController{
		CreatePostUsecase: usecase.NewCreatePostUsecase(pr, cr, timeout),
		Controller:        ctrl,
	}

	mux.HandleFunc(cpc.Data.Endpoints.CreatePostEndpoint, cpc.UserIdentity(cpc.Authorized(cpc.CreatePostController)))
}
