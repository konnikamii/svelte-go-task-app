package middleware

import "context"

type contextKey string

const userIDKey contextKey = "userID"

// ContextWithUserID stores the authenticated user ID in the request context.
func ContextWithUserID(ctx context.Context, userID int32) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext retrieves the authenticated user ID. Returns 0 if not set.
func UserIDFromContext(ctx context.Context) int32 {
	if id, ok := ctx.Value(userIDKey).(int32); ok {
		return id
	}
	return 0
}
