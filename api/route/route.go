package route

import (
	"database/sql"
	"net/http"
	"time"

	"forum/api/controller"
	"forum/bootstrap"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	NewCreatePostRouter(env, timeout, db, mux, ctrl)
	NewLoginRouter(env, timeout, db, mux, ctrl)
	NewLogoutRouter(env, timeout, db, mux, ctrl)
	NewCommentCreateRouter(env, timeout, db, mux, ctrl)
	NewListPostsRouter(env, timeout, db, mux, ctrl)
	NewPostReactionRouter(env, timeout, db, mux, ctrl)
	NewPostRouter(env, timeout, db, mux, ctrl)
	NewSignupRouter(env, timeout, db, mux, ctrl)
	NewCommentReactionRouter(env, timeout, db, mux, ctrl)
	NewCssRouter(mux)
}
