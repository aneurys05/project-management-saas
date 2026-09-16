package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "user_id"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			claims := &Claims{}

			token, err := jwt.ParseWithClaims(
				tokenString,
				claims,
				func(token *jwt.Token) (any, error) {

					if token.Method != jwt.SigningMethodHS256 {
						return nil, errors.New("unexpected signing method")
					}

					return []byte(secret), nil
				},
			)

			if err != nil || !token.Valid {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				userIDKey,
				claims.UserID,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
