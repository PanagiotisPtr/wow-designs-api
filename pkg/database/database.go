package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// JWTKey used to sign JSON Web Tokens
var JWTKey = []byte("secret_key")

// Claims used to encode jwt -- probably need to store this code on another file responsible for jwt
type Claims struct {
	Email string `json:"username"`
	jwt.StandardClaims
}

// Store is our connection struct for interfacing with the database
type Store struct {
	*mongo.Database
	client *mongo.Client
}

// Close a connection with the database
func (s *Store) Close(ctx context.Context) {
	log.Println("Closing connection")
	s.client.Disconnect(ctx)
}

// New connection with the database
func New(uri string, name string) (*Store, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to the database")
		return nil, err
	}

	db := client.Database(name)
	return &Store{db, client}, nil
}

// UserDetails that can be retrieved by the api
type UserDetails struct {
	Email       string
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth string
	SendDeals   bool
}

// SessionToken is just a Jason Web Token for a user session
type SessionToken struct {
	Token string
}

// UserCredentials include email and password for a user
type UserCredentials struct {
	Email    string
	Password string
}

// User shape
type User struct {
	Email       string
	Password    string
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth string
	SendDeals   bool
}

func userDetailsFromUser(user User) UserDetails {
	details := UserDetails{
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Gender:      user.Gender,
		DateOfBirth: user.DateOfBirth,
		SendDeals:   user.SendDeals,
	}

	return details
}

// GetSessionToken returns a token for a given session
func (s *Store) GetSessionToken(email string, password string) (SessionToken, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		return SessionToken{}, err
	}

	return SessionToken{Token: tokenString}, nil
}

func (s *Store) saveUserDetails(details UserDetails) error {
	userCollection := s.Collection("userDetails")

	insertResult, err := userCollection.InsertOne(context.TODO(), details)
	if err != nil {
		return err
	}

	log.Println("Added user details with ID: ", insertResult.InsertedID)

	return nil
}

func (s *Store) saveUserCredentials(creds UserCredentials) error {
	userCollection := s.Collection("userCredentials")

	insertResult, err := userCollection.InsertOne(context.TODO(), creds)
	if err != nil {
		return err
	}

	log.Println("Added user credentials with ID: ", insertResult.InsertedID)

	return nil
}

// CreateUser given login credentials and details
func (s *Store) CreateUser(creds UserCredentials, details UserDetails) error {
	err := s.saveUserCredentials(creds)
	if err != nil {
		return err
	}

	err = s.saveUserDetails(details)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser from the databse with email
func (s *Store) DeleteUser(email string) error {

	filter := bson.M{
		"email": bson.M{
			"$eq": email,
		},
	}

	userCollection := s.Collection("userDetails")

	deleteResult, err := userCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted %v documents in the userDetails collection\n", deleteResult.DeletedCount)

	userCollection = s.Collection("userCredentials")

	deleteResult, err = userCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted %v documents in the userCredentials collection\n", deleteResult.DeletedCount)

	return nil
}

// ChangeUserPassword for a specific user
func (s *Store) ChangeUserPassword(email string, newPassword string) error {
	filter := bson.M{
		"email": bson.M{
			"$eq": email,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"password": newPassword,
		},
	}

	userCollection := s.Collection("userCredentials")

	updateResult, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	log.Printf("Changed password for user %s. %v documents were matched and %v were modified.\n",
		email, updateResult.MatchedCount, updateResult.ModifiedCount)

	return nil
}

// ChangeUserDetails updates the details of a user given their email address
func (s *Store) ChangeUserDetails(email string, details UserDetails) error {
	filter := bson.M{
		"email": bson.M{
			"$eq": email,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"email":       details.Email,
			"firstname":   details.FirstName,
			"lastname":    details.LastName,
			"gender":      details.Gender,
			"dateofbirth": details.DateOfBirth,
			"senddeals":   details.SendDeals,
		},
	}

	userCollection := s.Collection("userDetails")

	updateResult, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return nil
}

func (s *Store) ChangeUserEmail(email string, newEmail string) error {
	filter := bson.M{
		"email": bson.M{
			"$eq": email,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"email": newEmail,
		},
	}

	userCollection := s.Collection("userDetails")

	updateResult, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	userCollection = s.Collection("userCredentials")

	updateResult, err = userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return nil
}

// GetUserPassword sends back the (hashed) password of a particular user
func (s *Store) GetUserPassword(email string) (UserCredentials, error) {
	userCollection := s.Collection("userCredentials")

	var result UserCredentials

	err := userCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&result)
	if err != nil {
		return UserCredentials{}, err
	}

	return result, nil
}

// GetUserDetails retuns the user emails based on their email
func (s *Store) GetUserDetails(email string) (UserDetails, error) {
	userCollection := s.Collection("userDetails")

	var result User

	err := userCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&result)
	if err != nil {
		return UserDetails{}, err
	}

	return userDetailsFromUser(result), nil
}
