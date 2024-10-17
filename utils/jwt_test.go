package utils

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	userID := int64(123)
	email := "satrio@gmail.com"
	secret := "test_secret"
	expiryHours := 24

	token, err := GenerateJWT(userID, email, secret, expiryHours)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		assert.Equal(t, userID, int64(claims["UserID"].(float64)))
		assert.Equal(t, email, claims["Email"].(string))
		assert.Equal(t, time.Now().Add(time.Duration(expiryHours)*time.Hour).Unix(), int64(claims["exp"].(float64)))
	} else {
		t.Fatal("Claims are not valid")
	}
}
