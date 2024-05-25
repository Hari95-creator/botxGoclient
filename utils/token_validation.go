// utils/token_validation.go

package utils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	UserGid      string    `json:"user_gid"`
	CreationDate time.Time `json:"creation_date"`
	jwt.RegisteredClaims
}

func IsTokenValid(tokenString string) error {
	claims := &Claims{}

	// Parse the JWT token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return GetClientPublicKey(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return httpError("Invalid token signature", http.StatusUnauthorized)
		}
		return httpError("Failed to parse token", http.StatusBadRequest)
	}
	if !token.Valid {
		return httpError("Invalid token", http.StatusUnauthorized)
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return httpError("Token has expired", http.StatusUnauthorized)
	}

	if claims.CreationDate.IsZero() || time.Now().Before(claims.CreationDate) {
		return httpError("Invalid token creation date", http.StatusUnauthorized)
	}
	return nil
}
func httpError(message string, statusCode int) error {
	return &httpErrorString{message, statusCode}
}

type httpErrorString struct {
	message    string
	statusCode int
}

func (e *httpErrorString) Error() string {
	return e.message
}
func (e *httpErrorString) StatusCode() int {
	return e.statusCode
}
