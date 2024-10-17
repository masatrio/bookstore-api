package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID int64
	Email  string
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for the given user ID and email.
func GenerateJWT(userID int64, email, secret string, expiryHours int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expiryHours) * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
