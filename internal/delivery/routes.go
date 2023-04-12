package delivery

import "net/http"

func (h *Handler) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.userIdentity(h.Home))
	mux.HandleFunc("/auth/signup", h.Signup)
	mux.HandleFunc("/auth/login", h.Login)
	mux.HandleFunc("/auth/logout", h.userIdentity(h.authorized(h.Logout)))

	mux.HandleFunc("/posts/create", h.userIdentity(h.authorized(h.CreatePostHandler)))
	mux.HandleFunc("/posts", h.ListPostsHandler)
	mux.HandleFunc("/posts/", h.userIdentity(h.PostHandler))
	mux.HandleFunc("/post/like/", h.LikePost)
	mux.HandleFunc("/post/dislike/", h.userIdentity(h.authorized(h.DislikePost)))

	//comment
	//likes
	//
}
