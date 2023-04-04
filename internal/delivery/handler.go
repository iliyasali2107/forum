package delivery

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/validator"
	"html/template"
	"net/http"
	"sync"

	"forum/internal/service"

	"forum/pkg/logger"
)

type Handler struct {
	tmpl      *template.Template
	Service   *service.Service
	logger    *logger.Logger
	wg        sync.WaitGroup
	cfg       config
	validator *validator.Validator
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		tmpl:    template.Must(template.ParseGlob("")),
		Service: service,
	}
}

func (h *Handler) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/signup", h.Signup)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		h.logger.PrintError(fmt.Errorf("%w", "not found"), nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.logger.PrintInfo("Signup:GET", nil)
		return
	case http.MethodPost:
		user := &models.User{}
		user.Name = r.FormValue("name")
		user.Email = r.FormValue("email")
		user.Password.Plaintext = r.FormValue("password")

		errs := h.Service.AuthService.Signup(h.validator, user)

		if errs[0] == service.ErrUserExists {
			h.logger.PrintError(fmt.Errorf("%s", http.StatusText(http.StatusConflict)), nil)
			return
		}

		if errs[len(errs)-1] == service.ErrInternalServer {
			h.logger.PrintError(fmt.Errorf("%s", http.StatusText(http.StatusInternalServerError)), nil)
			return
		}

		if len(errs) > 1 {
			//TODO: should render form with each error on its own field
			h.logger.PrintError(fmt.Errorf("%s", "signup: invalid form"), nil)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		h.logger.PrintError(fmt.Errorf("%s", "signup: method not allowed"), nil)
		return
	}
}
