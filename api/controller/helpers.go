package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/domain/models"
)

var (
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidUsername     = errors.New("invalid username")
	ErrInvalidUsernameLen  = errors.New("username length out of range 32")
	ErrInvalidUsernameChar = errors.New("invalid username characters")
	ErrInternalServer      = errors.New("internal server error")
	ErrConfirmPassword     = errors.New("password doesn't match")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserExists          = errors.New("user already exists")
	ErrFormValidation      = errors.New("form validation failed")
)

// Context helpers
const ctxKeyUser ctxKey = "user"

type ctxKey string

func (h *Controller) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), ctxKeyUser, user)
	return r.WithContext(ctx)
}

func (h *Controller) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(ctxKeyUser).(*models.User)
	if !ok {
		return nil
	}

	return user
}

// Error Response helpers
func (lc *Controller) errorPage(w http.ResponseWriter, code int) {
	w.WriteHeader(code)

	data := struct {
		Status  int
		Message string
	}{
		Status:  code,
		Message: http.StatusText(code),
	}

	if err := lc.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		lc.logger.PrintError(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (lc *Controller) logError(err error) {
	lc.logger.PrintError(err)
}

func (lc *Controller) errorResponse(w http.ResponseWriter, status int) {
	lc.errorPage(w, status)
}

func (lc *Controller) ResponseServerError(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusInternalServerError)
}

func (lc *Controller) ResponseNotFound(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusNotFound)
}

func (lc *Controller) ResponseMethodNotAllowed(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusMethodNotAllowed)
}

func (lc *Controller) ResponseBadRequest(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusBadRequest)
}

func (lc *Controller) ResponseFailedValidation(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusUnprocessableEntity)
}

func (lc *Controller) ResponseEditConflict(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusConflict)
}

func (lc *Controller) ResponseRateLimitExceeded(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusTooManyRequests)
}

func (lc *Controller) ResponseInvalidCredentials(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusUnauthorized)
}

func (lc *Controller) ResponseInvalidAuthenticationToken(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	lc.errorResponse(w, http.StatusUnauthorized)
}

func (lc *Controller) ResponseUnauthorized(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusUnauthorized)
}

func (lc *Controller) ResponseInactiveAccount(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusForbidden)
}

func (lc *Controller) ResponseNotPermitted(w http.ResponseWriter) {
	lc.errorResponse(w, http.StatusForbidden)
}

// render
func (h *Controller) render(w http.ResponseWriter, name string, td any) {
	err := h.tmpl.ExecuteTemplate(w, name, td)
	if err != nil {
		h.logger.PrintInfo("render: " + err.Error())
		h.ResponseServerError(w)
		return
	}
}

// url
func GetIdFromURL(path string) (int, error) {
	s := strings.Split(path, "/")

	if len(s) <= 3 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	if len(s[3:]) > 1 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	id, err := strconv.Atoi(s[3])
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetIdFromShortURL(path string) (int, error) {
	s := strings.Split(path, "/")
	if len(s) <= 2 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	if len(s[2:]) > 1 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	id, err := strconv.Atoi(s[2])
	if err != nil {
		return 0, err
	}

	return id, nil
}

// TODO: replace all GetIdFromURL and GetIdFromShortURL with function below
func GetIdFromURL2(numOfWords int, path string) (int, error) {
	s := strings.Split(path, "/")
	if len(s) <= numOfWords+1 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	if len(s[numOfWords+1:]) > 1 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	id, err := strconv.Atoi(s[numOfWords+1])
	if err != nil {
		return 0, err
	}

	return id, nil
}
