package delivery

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"forum/internal/models"
	"forum/pkg/validator"

	"github.com/felixge/httpsnoop"
	"golang.org/x/time/rate"
)

const ctxKeyUser ctxKey = "user"

type ctxKey string

func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		var err error
		c, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, models.User{})))
				return
			}
			h.errorPage(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err = h.Service.ParseToken(c.Value)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, models.User{})))
			return
		}
		if user.Expires.Before(time.Now()) {
			if err := h.Service.DeleteToken(c.Value); err != nil {
				h.errorPage(w, http.StatusInternalServerError, err.Error())
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, models.User{})))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	}
}

func (h *Handler) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				h.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				h.serverErrorResponse(w, r, err)
			}

			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(h.config.limiter.rps), h.config.limiter.burst)}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				h.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = h.contextSetUser(r, models.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			h.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.NewValidator()

		if ValidateTokenPlaintext(v, token); !v.Valid() {
			h.invalidCredentialsResponse(w, r)
			return
		}

		user, err := h.models.Users.GetForToken(ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				h.invalidCredentialsResponse(w, r)
			default:
				h.serverErrorResponse(w, r, err)
			}
			return
		}

		r = h.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := h.contextGetUser(r)

		if user.IsAnonymous() {
			h.authenticationRequireResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// func (h *Handler) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
// 	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user := h.contextGetUser(r)

// 		if !user.Activated {
// 			h.inactiveAccountResponse(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})

// 	return h.requireAuthenticatedUser(fn)
// }

// func (h *Handler) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		user := h.contextGetUser(r)

// 		permissions, err := h.models.Permissions.GetAllForUser(user.ID)
// 		if err != nil {
// 			h.serverErrorResponse(w, r, err)
// 			return
// 		}

// 		if !permissions.Include(code) {
// 			h.notPermittedResponse(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	}

// 	return h.requireActivatedUser(fn)
// }

func (h *Handler) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		origin := r.Header.Get("Origin")

		if origin != "" && len(h.config.cors.trustedOrigins) != 0 {
			for i := range h.config.cors.trustedOrigins {
				if origin == h.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_Ms")
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totalRequestsReceived.Add(1)

		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalResponsesSent.Add(1)

		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())

		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}
