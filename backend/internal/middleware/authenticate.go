package middleware

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/handlers"
)

// RequireAuth is an HTTP middleware that validates the session token from cookie
// against server-side session state in the database and
// injects the user ID into the request context.
// Apply it to any route group that requires a valid session.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bh := &handlers.BaseHandler{}

		token, err := SessionTokenFromRequest(r)
		if err != nil {
			bh.AppError(w, apperrors.Unauthenticated("missing or invalid session"))
			return
		}

		userID, err := SessionUserID(r.Context(), token)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				bh.AppError(w, apperrors.Unauthenticated("invalid or expired session"))
				return
			}

			bh.AppError(w, apperrors.Internal("could not validate session"))
			return
		}

		if userID == 0 {
			bh.AppError(w, apperrors.Unauthenticated("invalid or expired session"))
			return
		}

		ctx := ContextWithUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
