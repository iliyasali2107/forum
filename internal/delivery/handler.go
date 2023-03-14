package delivery

import (
	"html/template"
	"net/http"

	"forum/internal/service"
)

type Handler struct {
	tmpl    *template.Template
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		tmpl:    template.Must(template.ParseGlob("web/template/*.html")),
		Service: service,
	}
}

func (h *Handler) InitRoutes(mux *http.ServeMux) {
}
