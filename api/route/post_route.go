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

func NewPostRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	pr := repository.NewPostRepository(db)
	ur := repository.NewUserRepository(db)
	cr := repository.NewCategoryRepository(db)
	rr := repository.NewReactionRepository(db)
	cmr := repository.NewCommentRepository(db)

	pdc := controller.PostController{
		PostDetailsUsecase:     usecase.NewPostDetailsUsecase(pr, ur, cr, rr, cmr, timeout),
		CommentReactionUsecase: usecase.NewCommentReactionUsecase(rr, timeout),
		Controller:             ctrl,
	}

	mux.HandleFunc(ctrl.Data.Endpoints.PostDetailsEndpoint, pdc.UserIdentity(pdc.PostController))
}
