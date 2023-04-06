package delivery

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/service"
	"forum/pkg/validator"
	"html/template"
	"net/http"
	"os"
	"sync"

	"forum/pkg/logger"
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
		logger:    logger.NewLogger(os.Stdout, 0),
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		h.logger.PrintError("handler:home: " + err.Error())
		h.ResponseServerError(w)
	}
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
		fmt.Println("error: not found")
		h.ResponseNotFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := h.tmpl.ExecuteTemplate(w, "signup.html", nil)
		if err != nil {
			fmt.Println("handler:signup: " + err.Error())
			h.ResponseServerError(w)
			return
		}
	case http.MethodPost:
		user := &models.User{}
		user.Name = r.FormValue("name")
		user.Email = r.FormValue("email")
		user.Password.Plaintext = r.FormValue("password")

		errs := h.Service.AuthService.Signup(h.validator, user)
		fmt.Println(errs)
		if len(errs) > 0 && errs[0] == service.ErrUserExists {
			h.logger.PrintError(http.StatusText(http.StatusConflict))
			h.ResponseServerError(w)
			return
		}

		if len(errs) > 0 && errs[len(errs)-1] == service.ErrInternalServer {
			h.logger.PrintError(http.StatusText(http.StatusInternalServerError))
			h.ResponseServerError(w)
			return
		}

		if len(errs) > 1 {
			//TODO: should render form with each error on its own field !!!!!!!!!!!!!!!!!!!!!!!
			h.logger.PrintError("signup: invalid form")
			return
		}
		fmt.Println(user)
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		h.logger.PrintError("signup: method not allowed")
		h.ResponseMethodNotAllowed(w)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/login" {
		fmt.Println("handler.login: not found")
		h.ResponseNotFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := h.tmpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			h.ResponseServerError(w)
			return
		}
	case http.MethodPost:
		user := &models.User{}
		user.Email = r.FormValue("email")
		user.Password.Plaintext = r.FormValue("password")

		err := h.Service.AuthService.Login(h.validator, user)
		if err != nil {
			switch err {
			case service.ErrUserNotFound:
				fmt.Println("handler:login: user not found")
				h.ResponseBadRequest(w)
				return
			case service.ErrInvalidPassword:
				fmt.Println("handler:login: password is not correct")
				h.ResponseBadRequest(w)
				return
			default:
				fmt.Println(err)
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
		fmt.Println("handler:logout: not found")
		h.ResponseNotFound(w)
		return
	}

	c, err := r.Cookie("access_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			fmt.Println("handler:logout: unauthorized")
			h.errorPage(w, http.StatusUnauthorized)
			return
		}
		fmt.Println("handler:logout: " + err.Error())
		h.ResponseServerError(w)
		return

	}
	err = h.Service.AuthService.DeleteToken(c.Value)
	if err != nil {
		fmt.Println("handler:logout: couldn't delete token" + err.Error())
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
