package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user-api/pkg/database"
	"user-api/pkg/gql"

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
			return
		}

		// Authenticate user with JWT same function as above but not a handler
		// Get user email from JWT
		cookie, err := r.Cookie("token")
		if err != nil {
			// It's ok if we don't get a cookie. Some resolvers don't require authentication
			if err != http.ErrNoCookie {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
		}
		result := gql.ExecuteQuery(rBody.Query, *s.GqlSchema, cookie)

		render.JSON(w, r, result)
	}
}
