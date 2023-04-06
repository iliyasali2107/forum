package delivery

import (
	"net/http"
)

func (h *Handler) errorPage(w http.ResponseWriter, code int) {
	w.WriteHeader(code)

	data := struct {
		Status  int
		Message string
	}{
		Status:  code,
		Message: http.StatusText(code),
	}

	if err := h.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		h.logger.PrintError(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Handler) logError(err error) {
	h.logger.PrintError(err.Error())
}

func (h *Handler) errorResponse(w http.ResponseWriter, status int) {
	h.errorPage(w, status)
}

func (h *Handler) ResponseServerError(w http.ResponseWriter) {
	h.logger.PrintError(http.StatusText(http.StatusInternalServerError))
	h.errorResponse(w, http.StatusInternalServerError)
}

func (h *Handler) ResponseNotFound(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusNotFound)
}

func (h *Handler) ResponseMethodNotAllowed(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusMethodNotAllowed)
}

func (h *Handler) ResponseBadRequest(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusBadRequest)
}

func (h *Handler) ResponseFailedValidation(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusUnprocessableEntity)
}

func (h *Handler) ResponseEditConflict(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusConflict)
}

func (h *Handler) ResponseRateLimitExceeded(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusTooManyRequests)
}

func (h *Handler) ResponseInvalidCredentials(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusUnauthorized)
}

func (h *Handler) ResponseInvalidAuthenticationToken(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	h.errorResponse(w, http.StatusUnauthorized)
}

func (h *Handler) ResponseUnauthorizedRequire(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusUnauthorized)
}

func (h *Handler) ResponseInactiveAccount(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusForbidden)
}

func (h *Handler) ResponseNotPermitted(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusForbidden)
}
