package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/masatrio/bookstore-api/config"
)

type Claims struct {
	UserID int64
	Email  string
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for the given user ID and email.
func GenerateJWT(userID int64, email string) (string, error) {
	secretKey := config.LoadConfig().JWT.Secret
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
