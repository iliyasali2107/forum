package delivery

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/service"
	"forum/pkg/logger"
	"forum/pkg/validator"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Handler struct {
	tmpl      *template.Template
	Service   *service.Service
	logger    *logger.Logger
	cfg       config
	validator *validator.Validator
	wg        sync.WaitGroup
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		tmpl:      template.Must(template.ParseGlob("./templates/*")),
		Service:   service,
		validator: validator.NewValidator(),
		logger:    logger.NewLogger(os.Stdout, logger.LevelInfo),
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Service.PostService.GetAllPosts()
	if err != nil {
		h.logger.PrintError(err)
		h.ResponseServerError(w)
		return
	}
	err = h.tmpl.ExecuteTemplate(w, "show_posts.html", posts)
	if err != nil {
		h.logger.PrintError(err)
		h.ResponseServerError(w)
	}
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
		h.logger.PrintError(fmt.Errorf("handler: signup: not found"))
		h.ResponseNotFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := h.tmpl.ExecuteTemplate(w, "signup.html", nil)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: signup: internal server error"))
			h.ResponseServerError(w)
			return
		}
	case http.MethodPost:
		user := &models.User{}
		user.Name = r.FormValue("name")
		user.Email = r.FormValue("email")
		user.Password.Plaintext = r.FormValue("password")

		err := h.Service.AuthService.Signup(h.validator, user)
		if err == service.ErrUserExists {
			h.logger.PrintError(err)
			h.ResponseEditConflict(w)
			return
		}

		if err == service.ErrInternalServer {
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}

		if err == service.ErrFormValidation {
			h.ResponseBadRequest(w)
			h.render(w, "signup.html", h.validator)
			h.validator.Errors = map[string]string{}
			return
		}

		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	default:
		h.logger.PrintError(fmt.Errorf("handler: signup: method not allowed"))
		h.ResponseMethodNotAllowed(w)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/login" {
		h.logger.PrintError(fmt.Errorf("handler: login: not found"))
		h.ResponseNotFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := h.tmpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}
	case http.MethodPost:
		user := &models.User{}
		user.Email = r.FormValue("email")
		user.Password.Plaintext = r.FormValue("password")

		err := h.Service.AuthService.Login(user)
		if err != nil {
			switch err {
			case service.ErrUserNotFound:
				h.logger.PrintError(fmt.Errorf("handler: login: user not found"))
				h.ResponseBadRequest(w)
				return
			case service.ErrInvalidPassword:
				h.logger.PrintError(fmt.Errorf("handler: login: password is not correct"))
				h.ResponseBadRequest(w)
				return
			default:
				h.logger.PrintError(fmt.Errorf("handler: login: password is not correct"))
				h.ResponseServerError(w)
				return
			}
		}

		cookie := http.Cookie{}
		cookie.Name = "access_token"
		cookie.Value = *user.Token
		cookie.Expires = *user.Expires
		cookie.Path = "/"
		cookie.HttpOnly = true
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/logout" {
		h.logger.PrintError(fmt.Errorf("handler: logout: not found"))
		h.ResponseNotFound(w)
		return
	}

	c, err := r.Cookie("access_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.logger.PrintError(fmt.Errorf("handler: logout: unauthorized"))
			h.errorPage(w, http.StatusUnauthorized)
			return
		}
		fmt.Println("handler: logout: " + err.Error())
		h.ResponseServerError(w)
		return
	}

	err = h.Service.AuthService.DeleteToken(c.Value)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: logout: couldn't delete token"))
		h.ResponseServerError(w)
		return
	}

	cookie := &http.Cookie{
		Name:   "access_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// PostsHandler TODO: invalid field messages for each field
func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/posts/create" {
		h.logger.PrintError(fmt.Errorf("handler: post-create: not found"))
		h.ResponseNotFound(w)
		return
	}

	user := h.contextGetUser(r)

	switch r.Method {
	case http.MethodGet:
		categories, err := h.Service.PostService.GetAllCategories()

		data := struct {
			Errors     map[string]string
			Categories []*models.Category
		}{
			Categories: categories,
		}

		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: post-create: GetAllCategories error"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}

		err = h.tmpl.ExecuteTemplate(w, "create_post.html", data)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: post-create: ExecuteTemplate error"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: post-create: ParseForm error"))
			h.logger.PrintError(err)
			h.ResponseBadRequest(w)
			return
		}

		post := &models.Post{}
		post.Categories = r.Form["category"]
		post.Title = r.FormValue("title")
		post.Content = r.FormValue("content")
		post.User = user

		err = h.Service.PostService.CreatePost(h.validator, post)
		if err != nil {
			if err == service.ErrFormValidation {
				categories, err := h.Service.PostService.GetAllCategories()
				if err != nil {
					h.logger.PrintError(fmt.Errorf("handler: post-create: CreatePost ErrFormValidation error"))
					h.logger.PrintError(err)
					h.ResponseServerError(w)
					return
				}

				data := struct {
					Errors     map[string]string
					Categories []*models.Category
				}{
					Errors:     h.validator.Errors,
					Categories: categories,
				}
				h.ResponseBadRequest(w)
				h.render(w, "create_post.html", data)
				return
			}

			// TODO: ------------------
			h.logger.PrintError(err)
			return
		}
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	default:
		h.logger.PrintError(fmt.Errorf("create-post: method not allowed"))
		h.ResponseMethodNotAllowed(w)
		return
	}
}

func (h *Handler) ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.PrintError(fmt.Errorf("handler: listPost: method not allowed"))
		h.ResponseMethodNotAllowed(w)
		return
	}

	if r.URL.Path != "/posts" {
		h.logger.PrintError(fmt.Errorf("handler: listPost: not found"))
		h.ResponseNotFound(w)
		return
	}

	user := h.contextGetUser(r)

	queryMap := r.URL.Query()
	filter, ok := queryMap["filter"]
	if !ok {
		posts, err := h.Service.PostService.GetAllPosts()
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: listPost: GetAllPosts error"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}
		err = h.tmpl.ExecuteTemplate(w, "show_posts.html", posts)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: listPost: ExecuteTemplate error"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}
		return
	}

	if user == nil {
		h.logger.PrintError(fmt.Errorf("handler: listPost: unauthorized"))
		h.ResponseUnauthorized(w)
		return
	}

	switch filter[0] {
	case "created":
		posts, err := h.Service.PostService.GetCreatedPosts(user.ID)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: listPost: created GetCreatedPosts"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}

		err = h.tmpl.ExecuteTemplate(w, "show_posts.html", posts)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: listPost: created ExecuteTemplate show_posts.html"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}
	case "liked":
		posts, err := h.Service.PostService.GetLikedPosts(user.ID)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: listPost: liked GetLikedPosts"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}

		err = h.tmpl.ExecuteTemplate(w, "show_posts.html", posts)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: listPost: liked ExecuteTemplate show_posts.html"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}
	case "disliked":
		posts, err := h.Service.PostService.GetDislikedPosts(user.ID)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: listPost: disliked GetLikedPosts"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}

		err = h.tmpl.ExecuteTemplate(w, "show_posts.html", posts)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: listPost: disliked ExecuteTemplate show_posts.html"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}
	default:
		h.logger.PrintError(fmt.Errorf("handler: listPost: bad request"))
		h.ResponseBadRequest(w)
		return
	}
}

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromShortURL(r.URL.Path)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: PostHandler: not found"))
		h.logger.PrintError(err)
		h.ResponseNotFound(w)
		return
	}

	post, err := h.Service.PostService.GetPost(id)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: PostHandler: GetPost"))
		h.logger.PrintError(err)
		h.ResponseNotFound(w)
		return
	}

	likes, err := h.Service.ReactionService.GetPostLikes(id)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: PostHandler: GetPostLikes"))
		h.logger.PrintError(err)
		h.ResponseServerError(w)
		return
	}

	dislikes, err := h.Service.ReactionService.GetPostDislikes(id)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: PostHandler: GetPostDislikes"))
		h.logger.PrintError(err)
		h.ResponseServerError(w)
		return
	}

	post.Likes = likes
	post.Dislikes = dislikes

	// comments, err := h.Service.CommentService.GetCommentsByPostId(id)
	// fmt.Println(comments)
	// if err != nil {
	// 	h.logger.PrintError(fmt.Errorf("handler: PostHandler: GetCommentsByPostId"))
	// 	h.logger.PrintError(err)
	// 	h.ResponseServerError(w)
	// 	return
	// }

	// post.Comments = comments

	err = h.tmpl.ExecuteTemplate(w, "post.html", post)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: PostHandler: ExecuteTemplate post.html"))
		h.logger.PrintError(err)
		h.ResponseServerError(w)
		return
	}
}

func (h *Handler) LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.PrintError(fmt.Errorf("handler: LikePost: method not allowed"))
		h.ResponseMethodNotAllowed(w)
		return
	}

	id, err := GetIdFromURL(r.URL.Path)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: LikePost: GetIdFromURL"))
		h.logger.PrintError(err)
		h.ResponseNotFound(w)
		return
	}

	post, err := h.Service.PostService.GetPost(id)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: LikePost: GetPost"))
		h.logger.PrintError(err)
		h.ResponseServerError(w)
		return
	}
	user := h.contextGetUser(r)

	reaction := &models.Reaction{User: user, Post: post, Type: models.LikeType}

	err = h.Service.ReactionService.LikePost(reaction)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: LikePost: h.Service.ReactionService.LikePost()"))
		h.logger.PrintError(err)
		h.ResponseServerError(w)
		return
	}

	http.Redirect(w, r, "/posts/"+strconv.Itoa(post.ID), http.StatusSeeOther)
}

func (h *Handler) DislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.PrintError(fmt.Errorf("handler: DislikePost: method not allowed"))
		h.ResponseMethodNotAllowed(w)
		return
	}

	id, err := GetIdFromURL(r.URL.Path)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: DislikePost: GetIdFromURL"))
		h.logger.PrintError(err)
		h.ResponseNotFound(w)
		return
	}
	post, err := h.Service.PostService.GetPost(id)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: DislikePost: GetPost"))
		h.logger.PrintError(err)
		h.ResponseServerError(w)
		return
	}
	user := h.contextGetUser(r)

	reaction := &models.Reaction{User: user, Post: post, Type: models.DislikeType}

	err = h.Service.ReactionService.DislikePost(reaction)
	if err != nil {
		h.logger.PrintError(fmt.Errorf("handler: DislikePost: h.Service.ReactionService.DislikePost()"))
		h.logger.PrintError(err)
		h.ResponseServerError(w)
		return
	}

	http.Redirect(w, r, "/posts/"+strconv.Itoa(post.ID), http.StatusSeeOther)
}

func (h *Handler) CommentHandler(w http.ResponseWriter, r *http.Request) {
	// id, err := GetIdFromShortURL(r.URL.Path)
	// if err != nil {
	// 	h.logger.PrintError(err)
	// 	h.ResponseNotFound(w)
	// 	return
	// }

	// post, err := h.Service.PostService.GetPost(id)
	// if err != nil {
	// 	h.logger.PrintError(err)
	// 	h.ResponseServerError(w)
	// 	return
	// }
	// user := h.contextGetUser(r)

	// comment, err := h.Service.CommentService.
}

func (h *Handler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		id, err := GetIdFromURL2(3, r.URL.Path)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: comment-create: GetIdFromURL2"))
			h.logger.PrintError(err)
			h.ResponseNotFound(w)
			return
		}

		post, err := h.Service.PostService.GetPost(id)
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: comment-create: GetPost"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}

		user := h.contextGetUser(r)

		err = r.ParseForm()
		if err != nil {
			h.logger.PrintError(fmt.Errorf("handler: comment-create: ParseForm"))
			h.logger.PrintError(err)
			h.ResponseBadRequest(w)
			return
		}

		comment := &models.Comment{}
		comment.Content = r.FormValue("content")
		comment.User = user
		comment.Post = post
		// comment.Parent

		err = h.Service.CommentService.CreateComment(h.validator, comment)
		if err != nil {
			if err == service.ErrFormValidation {
				h.logger.PrintError(fmt.Errorf("handler: comment-create: CreateComment ErrFormValidation"))
				h.logger.PrintError(err)
				h.ResponseBadRequest(w)
				return
			}
			h.logger.PrintError(fmt.Errorf("handler: comment-create: CreateComment"))
			h.logger.PrintError(err)
			h.ResponseServerError(w)
			return
		}

		http.Redirect(w, r, "/posts/"+strconv.Itoa(post.ID), http.StatusSeeOther)
	default:
		h.logger.PrintError(fmt.Errorf("handler: comment-create: method not allowed"))
		h.ResponseMethodNotAllowed(w)
		return
	}
}
