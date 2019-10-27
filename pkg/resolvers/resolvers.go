package resolvers

import (
	"fmt"
	"net/http"
	"user-api/pkg/database"

	"github.com/dgrijalva/jwt-go"
)

// Resolver struct wraps around a store object which is used to interface
// with the databse in order to get thew data
type Resolver struct {
	Store *database.Store
}

func getUserEmailFromCookie(cookie *http.Cookie) (string, error) {
	tokenString := cookie.Value
	claims := &database.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return database.JWTKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", fmt.Errorf("Invalid signature")
		}
		return "", fmt.Errorf("Bad request")
	}
	if !token.Valid {
		return "", fmt.Errorf("Invalid Token")
	}

	return claims.Email, nil
}
