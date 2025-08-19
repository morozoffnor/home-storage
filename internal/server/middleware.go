package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func homeCtx(next http.Handler) http.Handler {
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
