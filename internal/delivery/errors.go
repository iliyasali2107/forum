package delivery

import (
	"fmt"
	"net/http"
	"strings"
)

func (h *Handler) errorPage(w http.ResponseWriter, code int, errorText string) {
	w.WriteHeader(code)
	data := struct {
		Status  int
		Message string
		ErrText string
	}{
		Status:  code,
		Message: http.StatusText(code),
		ErrText: errorText,
	}
	if data.Status != http.StatusInternalServerError {
		temp := strings.Split(errorText, ":")
		data.ErrText = temp[len(temp)-1]
	}
	if err := h.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// fmt.Fprintf(w, "%d - %s\n", data.Status, data.Message)
	}
}

func (h *Handler) logError(r *http.Request, err error) {
	h.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (h *Handler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := h.writeJSON(w, status, env, nil)
	if err != nil {
		h.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	h.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (h *Handler) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	h.errorResponse(w, r, http.StatusNotFound, message)
}

func (h *Handler) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	h.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (h *Handler) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *Handler) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	h.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (h *Handler) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	h.errorResponse(w, r, http.StatusConflict, message)
}

func (h *Handler) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	h.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (h *Handler) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	h.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (h *Handler) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	h.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (h *Handler) authenticationRequireResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	h.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (h *Handler) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	h.errorResponse(w, r, http.StatusForbidden, message)
}

func (h *Handler) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	h.errorResponse(w, r, http.StatusForbidden, message)
}
