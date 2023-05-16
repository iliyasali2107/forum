package route

import (
	"database/sql"
	"net/http"
	"time"

	"forum/api/controller"
	"forum/bootstrap"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	NewCommentDetailsRouter(env, timeout, db, mux)
	NewCreatePostRouter(env, timeout, db, mux)
	NewListPostsRouter(env, timeout, db, mux)
	NewLoginRouter(env, timeout, db, mux)
	NewLogoutRouter(env, timeout, db, mux)
	NewCommentCreateRouter(env, timeout, db, mux)
	NewPostReactionRouter(env, timeout, db, mux)
	NewPostRouter(env, timeout, db, mux)
	NewSignupRouter(env, timeout, db, mux, ctrl)
}
