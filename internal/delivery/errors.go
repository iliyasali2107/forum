package delivery

import (
	"net/http"
)

func (h *Controller) errorPage(w http.ResponseWriter, code int) {
	w.WriteHeader(code)

	data := struct {
		Status  int
		Message string
	}{
		Status:  code,
		Message: http.StatusText(code),
	}

	if err := h.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		h.logger.PrintError(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Controller) logError(err error) {
	h.logger.PrintError(err)
}

func (h *Controller) errorResponse(w http.ResponseWriter, status int) {
	h.errorPage(w, status)
}

func (h *Controller) ResponseServerError(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusInternalServerError)
}

func (h *Controller) ResponseNotFound(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusNotFound)
}

func (h *Controller) ResponseMethodNotAllowed(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusMethodNotAllowed)
}

func (h *Controller) ResponseBadRequest(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusBadRequest)
}

func (h *Controller) ResponseFailedValidation(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusUnprocessableEntity)
}

func (h *Controller) ResponseEditConflict(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusConflict)
}

func (h *Controller) ResponseRateLimitExceeded(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusTooManyRequests)
}

func (h *Controller) ResponseInvalidCredentials(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusUnauthorized)
}

func (h *Controller) ResponseInvalidAuthenticationToken(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	h.errorResponse(w, http.StatusUnauthorized)
}

func (h *Controller) ResponseUnauthorized(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusUnauthorized)
}

func (h *Controller) ResponseInactiveAccount(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusForbidden)
}

func (h *Controller) ResponseNotPermitted(w http.ResponseWriter) {
	h.errorResponse(w, http.StatusForbidden)
}
