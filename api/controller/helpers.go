package controller

import (
	"context"
	"fmt"
	"forum/domain/models"
	"net/http"
	"strconv"
	"strings"
)

// Context helpers
const CtxKeyUser CtxKey = "user"

type CtxKey string

func (c *Controller) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), CtxKeyUser, user)
	return r.WithContext(ctx)
}

func (c *Controller) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(CtxKeyUser).(*models.User)
	if !ok {
		return nil
	}

	return user
}

// Error Response helpers
func (c *Controller) errorPage(w http.ResponseWriter, code int) {
	w.WriteHeader(code)

	data := struct {
		Status  int
		Message string
	}{
		Status:  code,
		Message: http.StatusText(code),
	}

	if err := c.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		c.logger.PrintError(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (c *Controller) logError(err error) {
	c.logger.PrintError(err)
}

func (c *Controller) errorResponse(w http.ResponseWriter, status int) {
	c.errorPage(w, status)
}

func (c *Controller) ResponseServerError(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusInternalServerError)
}

func (c *Controller) ResponseForbidden(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusForbidden)
}

func (c *Controller) ResponseNotFound(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusNotFound)
}

func (c *Controller) ResponseMethodNotAllowed(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusMethodNotAllowed)
}

func (c *Controller) ResponseBadRequest(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusBadRequest)
}

func (c *Controller) ResponseFailedValidation(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusUnprocessableEntity)
}

func (c *Controller) ResponseEditConflict(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusConflict)
}

func (c *Controller) ResponseRateLimitExceeded(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusTooManyRequests)
}

func (c *Controller) ResponseInvalidCredentials(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusUnauthorized)
}

func (c *Controller) ResponseInvalidAuthenticationToken(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	c.errorResponse(w, http.StatusUnauthorized)
}

func (c *Controller) ResponseUnauthorized(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusUnauthorized)
}

func (c *Controller) ResponseInactiveAccount(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusForbidden)
}

func (c *Controller) ResponseNotPermitted(w http.ResponseWriter) {
	c.errorResponse(w, http.StatusForbidden)
}

// render
func (c *Controller) render(w http.ResponseWriter, name string, td any) {
	err := c.tmpl.ExecuteTemplate(w, name, td)
	if err != nil {
		c.logger.PrintInfo("render: " + err.Error())
		c.ResponseServerError(w)
		return
	}
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
