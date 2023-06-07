package route

import (
	"database/sql"
	"forum/api/controller"
	"forum/bootstrap"
	"net/http"
	"time"
)

func NewCssRouter(env *bootstrap.Env, timeout time.Duration, db *sql.DB, mux *http.ServeMux, ctrl *controller.Controller) {
	mux.Handle("/ui/assets/", http.StripPrefix("/ui/assets/", http.FileServer(http.Dir("ui/assets"))))
}
