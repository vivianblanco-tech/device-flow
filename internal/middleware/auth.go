package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"os"

	"github.com/yourusername/laptop-tracking-system/internal/auth"
	"github.com/yourusername/laptop-tracking-system/internal/models"
)

// isProduction checks if the application is running in production
func isProduction() bool {
	env := os.Getenv("APP_ENV")
	return env == "production"
}

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserContextKey is the key for storing user in request context
	UserContextKey ContextKey = "user"
	// SessionContextKey is the key for storing session in request context
	SessionContextKey ContextKey = "session"
)

// SessionCookieName is the name of the session cookie
const SessionCookieName = "session_token"

// AuthMiddleware validates the session and adds user info to the request context
func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get session token from cookie
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				// No session cookie, continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			// Validate session
			session, err := auth.ValidateSession(r.Context(), db, cookie.Value)
			if err != nil {
				// Log error but don't block request
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

		if session == nil {
			// Invalid or expired session, clear cookie
			http.SetCookie(w, &http.Cookie{
				Name:     SessionCookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   isProduction(),
				SameSite: http.SameSiteStrictMode,
			})
			next.ServeHTTP(w, r)
			return
		}

			// Add session and user to context
			ctx := context.WithValue(r.Context(), SessionContextKey, session)
			ctx = context.WithValue(ctx, UserContextKey, session.User)

			// Continue with authenticated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth middleware ensures the user is authenticated
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireRole middleware ensures the user has one of the specified roles
func RequireRole(roles ...models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r.Context())
			if user == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Check if user has one of the required roles
			hasRole := false
			for _, role := range roles {
				if user.HasRole(role) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext retrieves the authenticated user from the request context
func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

// GetSessionFromContext retrieves the session from the request context
func GetSessionFromContext(ctx context.Context) *models.Session {
	session, ok := ctx.Value(SessionContextKey).(*models.Session)
	if !ok {
		return nil
	}
	return session
}

// IsAuthenticated checks if the request has an authenticated user
func IsAuthenticated(r *http.Request) bool {
	return GetUserFromContext(r.Context()) != nil
}

