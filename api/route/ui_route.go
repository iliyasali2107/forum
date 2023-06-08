package route

import (
	"net/http"
)

func NewCssRouter(mux *http.ServeMux) {
	mux.Handle("/ui/assets/", http.StripPrefix("/ui/assets/", http.FileServer(http.Dir("ui/assets"))))
}
