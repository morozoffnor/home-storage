package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/morozoffnor/home-storage/internal/auth"
	"github.com/morozoffnor/home-storage/internal/database"
)

type Middleware struct {
	auth *auth.Auth
	db   *database.Database
}

func NewMiddleware(a *auth.Auth, db *database.Database) *Middleware {
	return &Middleware{
		auth: a,
		db:   db,
	}
}

func (m *Middleware) homeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		homeIDStr := chi.URLParam(r, "home_id")
		var homeID int
		if _, err := fmt.Sscanf(homeIDStr, "%d", &homeID); err != nil {
			http.Error(w, "invalid home ID format", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "home_id", homeID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "user_id")
		var userID int
		if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
			http.Error(w, "invalid user id format", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("Authorization")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := m.auth.Jwt.ParseToken(cookie.Value)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// update token if it expires soon
		if claims.ExpiresAt.Add(time.Hour).After(time.Now()) {
			token, _ := m.auth.Jwt.GenerateToken(claims.UserEmail)
			ctx, _ := m.auth.Jwt.AddTokenToCookies(&w, r, token)
			r = r.WithContext(ctx)
		}

		ctx := context.WithValue(r.Context(), auth.ContextUserEmail, claims.UserEmail)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
