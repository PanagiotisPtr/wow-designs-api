package resolvers

import (
	"fmt"
	"log"
	"net/http"
	"user-api/pkg/database"

	"golang.org/x/crypto/bcrypt"

	"github.com/graphql-go/graphql"
)

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

// Register registers a new user and returns true if it was successful and false otherwise
func (r *Resolver) Register(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Args["email"].(string)

	if !ok {
		return false, fmt.Errorf("Resolver.Authentication: invalide resolve arguments: %v", p.Args)
	}

	password, ok := p.Args["password"].(string)

	if !ok {
		return false, fmt.Errorf("Resolver.Authentication: invalide resolve arguments: %v", p.Args)
	}

	_, err := r.Store.GetUserDetails(email)
	if err == nil {
		return false, fmt.Errorf("Resolver.Register: User already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return false, err
	}

	details, errors := getUserDetailsFromParams(p)
	if errors != nil {
		log.Println("Resolver.Register: Errors when getting user details: ")
		log.Println(errors)
	}

	creds := database.UserCredentials{Email: email, Password: string(hash)}

	err = r.Store.CreateUser(creds, details)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Authenticate authenticates a user using their email and password and returns a JWT
func (r *Resolver) Authenticate(p graphql.ResolveParams) (interface{}, error) {
	email, ok := p.Args["email"].(string)

	if !ok {
		return nil, fmt.Errorf("AuthenticationResolver: invalide resolve arguments: %v", p.Args)
	}

	password, ok := p.Args["password"].(string)

	if !ok {
		return nil, fmt.Errorf("AuthenticationResolver: invalide resolve arguments: %v", p.Args)
	}

	err := r.authenticateUser(email, password)
	if err != nil {
		return nil, err
	}

	token, err := r.Store.GetSessionToken(email, password)

	return token, err
}

// UserDetails sends back the details about a specific the user who sent the request
// based on the JWT that was included in the cookie
func (r *Resolver) UserDetails(p graphql.ResolveParams) (interface{}, error) {
	cookie, ok := p.Context.Value("cookie").(*http.Cookie)
	if !ok || cookie == nil {
		return nil, fmt.Errorf("Error parsing JWT cookie from header")
	}

	authEmail, err := getUserEmailFromCookie(cookie)
	if err != nil {
		return nil, err
	}

	userDetails, err := r.Store.GetUserDetails(authEmail)
	return userDetails, err
}

// ChangePassword changes the current user password. Returns true on success
// or false otherwise
func (r *Resolver) ChangePassword(p graphql.ResolveParams) (interface{}, error) {
	cookie := p.Context.Value("cookie").(*http.Cookie)
	authEmail, err := getUserEmailFromCookie(cookie)
	if err != nil {
		return nil, err
	}

	password, ok := p.Args["password"].(string)
	if !ok {
		return nil, fmt.Errorf("AuthenticationResolver: invalide resolve arguments: %v", p.Args)
	}

	newPassword, ok := p.Args["newPassword"].(string)
	if !ok {
		return nil, fmt.Errorf("AuthenticationResolver: invalide resolve arguments: %v", p.Args)
	}

	err = r.authenticateUser(authEmail, password)
	if err != nil {
		return false, fmt.Errorf("Resolver.ChangePassword: Could not authenticate user: %s", err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return false, fmt.Errorf("Resolver.ChangePassword: Dailed to hash new password: %s", err.Error())
	}

	err = r.Store.ChangeUserPassword(authEmail, string(hash))
	if err != nil {
		return false, fmt.Errorf("Resolver.ChangePassword: Could not change password: %s", err.Error())
	}

	return true, nil
}

// ChangeUserDetails changes the details of a specific user
func (r *Resolver) ChangeUserDetails(p graphql.ResolveParams) (interface{}, error) {
	cookie := p.Context.Value("cookie").(*http.Cookie)
	authEmail, err := getUserEmailFromCookie(cookie)
	if err != nil {
		return false, err
	}

	userDetails, err := r.Store.GetUserDetails(authEmail)
	if err != nil {
		return false, err
	}

	newUserDetails, _ := getUserDetailsFromParams(p)
	// Not allowed to change email without authentication
	newUserDetails.Email = userDetails.Email
	if newUserDetails.FirstName == "" {
		newUserDetails.FirstName = userDetails.FirstName
	}
	if newUserDetails.LastName == "" {
		newUserDetails.LastName = userDetails.LastName
	}
	if newUserDetails.Gender == "" {
		newUserDetails.Gender = userDetails.Gender
	}
	if newUserDetails.DateOfBirth == "" {
		newUserDetails.DateOfBirth = userDetails.DateOfBirth
	}
	log.Println(newUserDetails)
	err = r.Store.ChangeUserDetails(authEmail, newUserDetails)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Terminate a user account
func (r *Resolver) Terminate(p graphql.ResolveParams) (interface{}, error) {
	cookie := p.Context.Value("cookie").(*http.Cookie)
	authEmail, err := getUserEmailFromCookie(cookie)
	if err != nil {
		return nil, err
	}

	password, ok := p.Args["password"].(string)
	if !ok {
		return nil, fmt.Errorf("AuthenticationResolver: invalide resolve arguments: %v", p.Args)
	}

	err = r.authenticateUser(authEmail, password)
	if err != nil {
		return false, fmt.Errorf("Resolver.ChangePassword: Could not authenticate user: %s", err.Error())
	}

	err = r.Store.DeleteUser(authEmail)
	if err != nil {
		return false, fmt.Errorf("Resolver.Terminate: Could not delete user account: %s", err.Error())
	}

	return true, nil
}
