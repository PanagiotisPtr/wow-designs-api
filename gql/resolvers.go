package gql

import (
	"fmt"
	"log"
	"net/http"
	"webserver-init/database"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
)

type Resolver struct {
	store *database.Store
}

/*
func (r *Resolver) UserUpdateResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Args["email"].(string)
	if !ok {
		return nil, fmt.Errorf("Expected email parameter for resolver")

}*/

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

func getUserDetailsFromParams(p graphql.ResolveParams) (database.UserDetails, []error) {
	var details database.UserDetails
	var errors []error
	email, ok := p.Args["email"].(string)

	if !ok {
		err := fmt.Errorf("getUserDetailsFromParams: missing email in params")
		errors = append(errors, err)
	} else {
		details.Email = email
	}
	firstName, ok := p.Args["firstName"].(string)
	if !ok {
		err := fmt.Errorf("getUserDetailsFromParams: firstName email in params")
		errors = append(errors, err)
	} else {
		details.FirstName = firstName
	}
	lastName, ok := p.Args["lastName"].(string)

	if !ok {
		err := fmt.Errorf("getUserDetailsFromParams: missing lastName in params")
		errors = append(errors, err)
	} else {
		details.LastName = lastName
	}
	gender, ok := p.Args["gender"].(string)

	if !ok {
		err := fmt.Errorf("getUserDetailsFromParams: missing gender in params")
		errors = append(errors, err)
	} else {
		details.Gender = gender
	}
	dateOfBirth, ok := p.Args["dateOfBirth"].(string)

	if !ok {
		err := fmt.Errorf("getUserDetailsFromParams: missing dateOfBirth in params")
		errors = append(errors, err)
	} else {
		details.DateOfBirth = dateOfBirth
	}
	sendDeals, ok := p.Args["sendDeals"].(bool)

	if !ok {
		err := fmt.Errorf("getUserDetailsFromParams: missing sendDeals in params")
		errors = append(errors, err)
	} else {
		details.SendDeals = sendDeals
	}

	return details, errors
}

func (r *Resolver) RegisterResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Args["email"].(string)

	if !ok {
		return false, fmt.Errorf("AuthenticationResolver: invalide resolve arguments: %v", p.Args)
	}

	password, ok := p.Args["password"].(string)

	if !ok {
		return false, fmt.Errorf("AuthenticationResolver: invalide resolve arguments: %v", p.Args)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return false, err
	}

	details, errors := getUserDetailsFromParams(p)
	if errors != nil {
		log.Println("RegisterResolver: Errors when getting user details: ")
		log.Println(errors)
	}

	creds := database.UserCredentials{Email: email, Password: string(hash)}

	err = r.store.CreateUser(creds, details)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Resolver) AuthenticationResolver(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Args["email"].(string)

	if !ok {
		return nil, fmt.Errorf("AuthenticationResolver: invalide resolve arguments: %v", p.Args)
	}

	password, ok := p.Args["password"].(string)

	if !ok {
		return nil, fmt.Errorf("AuthenticationResolver: invalide resolve arguments: %v", p.Args)
	}

	creds, err := r.store.GetUserPassword(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	token, err := r.store.GetSessionToken(email, password)

	return token, err
}

func (r *Resolver) UserDetailsResolver(p graphql.ResolveParams) (interface{}, error) {
	cookie := p.Context.Value("cookie").(*http.Cookie)
	authEmail, err := getUserEmailFromCookie(cookie)
	if err != nil {
		return nil, err
	}

	users, err := r.store.GetUserDetailsByEmail(authEmail)
	return users, err
}

func (r *Resolver) ListResolver(p graphql.ResolveParams) (interface{}, error) {
	users, err := r.store.GetAllUserDetails()

	return users, err
}
