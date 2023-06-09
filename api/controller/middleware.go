package controller

import (
	"context"
	"errors"
	"fmt"
	"forum/domain/models"
	"log"
	"net/http"
	"time"
)

func (ctrl *Controller) UserIdentity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		var err error
		c, err := r.Cookie("access_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, nil)))
				return
			}
			ctrl.ResponseBadRequest(w)
			return
		}

		user, err = ctrl.TokenUsecase.ParseToken(c.Value)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, nil)))
			return
		}

		if user.Expires.Before(time.Now()) {
			if err := ctrl.TokenUsecase.DeleteToken(c.Value); err != nil {
				ctrl.ResponseServerError(w)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, nil)))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, user)))
	}
}

func (c *Controller) Authorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// user := r.Context().Value(ctxKeyUser).(*models.User)
		u := r.Context().Value(CtxKeyUser)

		if u == nil {
			fmt.Println("middleware:authorized: user is not authorized")
			c.ResponseUnauthorized(w)
			return
		}

		// user := u.(*models.User)

		next.ServeHTTP(w, r)
	}
}

func (c *Controller) NotAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// user := r.Context().Value(ctxKeyUser).(*models.User)
		u := r.Context().Value(CtxKeyUser)
		if u != nil {
			log.Println("middleware:authorized: user is authorized")
			c.ResponseForbidden(w)
			return
		}

		// user := u.(*models.User)

		next.ServeHTTP(w, r)
	}
}
