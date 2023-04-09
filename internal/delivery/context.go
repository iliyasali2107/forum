package delivery

import (
	"context"
	"net/http"

	"forum/internal/models"
)

const ctxKeyUser ctxKey = "user"

type ctxKey string

func (h *Handler) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), ctxKeyUser, user)
	return r.WithContext(ctx)
}

func (h *Handler) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(ctxKeyUser).(*models.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
