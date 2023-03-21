package delivery

import (
	"html/template"
	"net/http"
	"sync"

	"forum/internal/service"

	"forum/pkg/logger"
)

type Handler struct {
	tmpl    *template.Template
	Service *service.Service
	logger  *logger.Logger
	wg      sync.WaitGroup
	config  config
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		tmpl:    template.Must(template.ParseGlob("web/template/*.html")),
		Service: service,
	}
}

func (h *Handler) InitRoutes(mux *http.ServeMux) {
}
