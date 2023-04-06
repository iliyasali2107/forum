package delivery

import "net/http"

func (h *Handler) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.userIdentity(h.Home))
	mux.HandleFunc("/auth/signup", h.Signup)
	mux.HandleFunc("/auth/login", h.Login)
	mux.HandleFunc("/auth/logout", h.userIdentity(h.authorized(h.Logout)))

	//post
	//mux.HandleFunc("post")

	//comment
	//likes
	//
}
