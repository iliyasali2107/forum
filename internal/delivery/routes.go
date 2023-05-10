package delivery

import "net/http"

func (h *Handler) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.userIdentity(h.Home))
	mux.HandleFunc("/auth/signup", h.Signup)
	mux.HandleFunc("/auth/login", h.Login)
	mux.HandleFunc("/auth/logout", h.userIdentity(h.authorized(h.Logout)))

	// post
	mux.HandleFunc("/posts/create", h.userIdentity(h.authorized(h.CreatePostHandler)))
	mux.HandleFunc("/posts", h.userIdentity(h.authorized(h.ListPostsHandler)))
	mux.HandleFunc("/posts/", h.userIdentity(h.PostHandler))
	mux.HandleFunc("/posts/reaction", h.userIdentity(h.authorized(h.ReactionPostHandler)))

	mux.HandleFunc("/posts/comment", h.userIdentity(h.authorized(h.CreateCommentHandler)))
	mux.HandleFunc("/posts/comments/", h.userIdentity(h.authorized(h.ShowCommentHandler)))
}
