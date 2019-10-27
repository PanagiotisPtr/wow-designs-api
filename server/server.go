package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webserver-init/database"
	"webserver-init/gql"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
)

type Server struct {
	GqlSchema *graphql.Schema
}

type reqBody struct {
	Query string `json:"query"`
}

func (s *Server) AuthenticatedUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Test")
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("No cookie")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenString := c.Value
		claims := &database.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return database.JWTKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Println(claims.Email)
		w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Email)))
	}
}

func getUserEmailFromRequest(r *http.Request) (string, int32, error) {
	fmt.Println("Test")
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", http.StatusUnauthorized, fmt.Errorf("No cookie found")
		}

		return "", http.StatusBadRequest, fmt.Errorf("Bad request")
	}

	tokenString := c.Value
	claims := &database.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return database.JWTKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", http.StatusUnauthorized, fmt.Errorf("Invalid signature")
		}
		return "", http.StatusBadRequest, fmt.Errorf("Bad request")
	}
	if !token.Valid {
		return "", http.StatusUnauthorized, fmt.Errorf("Invalid Token")
	}
	return claims.Email, http.StatusOK, nil
}

func (s *Server) GraphQL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Must provide graphql query in request body", 400)
			return
		}

		var rBody reqBody

		err := json.NewDecoder(r.Body).Decode(&rBody)
		if err != nil {
			http.Error(w, "Error parsing JSON request body", 400)
		}

		// Authenticate user with JWT same function as above but not a handler
		// Get user email from JWT
		userEmail := "user@email.com"
		result := gql.ExecuteQuery(rBody.Query, *s.GqlSchema, userEmail)

		render.JSON(w, r, result)
	}
}
