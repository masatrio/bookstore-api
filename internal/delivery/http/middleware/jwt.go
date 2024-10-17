package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/masatrio/bookstore-api/config"
	"github.com/masatrio/bookstore-api/utils"
)

type contextKey string

const userIDKey contextKey = "userID"

// JWTMiddleware checks the validity of the JWT token in the Authorization header.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, "Bearer ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := bearerToken[1]

		secretKey := config.LoadConfig().JWT.Secret

		token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*utils.Claims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			r = r.WithContext(ctx)
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext retrieves the user ID from the context.
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}