package delivery

import (
	"context"
	"errors"
	"fmt"
	"forum/internal/models"
	"net/http"
	"time"
)

// TODO: signup and login for authorized user not available

func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		var err error
		c, err := r.Cookie("access_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, nil)))
				return
			}
			h.ResponseBadRequest(w)
			return
		}

		user, err = h.Service.ParseToken(c.Value)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, nil)))
			return
		}
		if user.Expires.Before(time.Now()) {
			if err := h.Service.DeleteToken(c.Value); err != nil {
				h.ResponseServerError(w)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, nil)))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	}
}

func (h *Handler) recoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				h.ResponseServerError(w)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) authorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// user := r.Context().Value(ctxKeyUser).(*models.User)
		u := r.Context().Value(ctxKeyUser)

		if u == nil {
			fmt.Println("middleware:authorized: user is not authorized")
			h.ResponseUnauthorized(w)
			return
		}

		// user := u.(*models.User)

		next.ServeHTTP(w, r)
	}
}
