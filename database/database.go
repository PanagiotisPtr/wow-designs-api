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

func (s *Store) Close(ctx context.Context) {
	log.Println("Closing connection")
	s.client.Disconnect(ctx)
}

func New(uri string, name string) (*Store, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to the database")
		return nil, err
	}

	db := client.Database(name)
	log.Println("Successfully connected to the database")
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

func (s *Store) GetSessionToken(email string, password string) (SessionToken, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
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

func (s *Store) DeleteUser(email string) error {
	userCollection := s.Collection("userDetails")

	deleteResult, err := userCollection.DeleteOne(context.TODO(), bson.D{{Key: "email", Value: email}})
	if err != nil {
		return err
	}

	fmt.Printf("Deleted %v documents in the users collection\n", deleteResult.DeletedCount)

	return nil
}

func (s *Store) AuthenticateUser(email string, password string) (bool, error) {
	userCollection := s.Collection("userDetails")

	var user User

	err := userCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&user)
	if err != nil {
		return false, err
	}

	if user.Password == password {
		return true, nil
	}

	return false, nil
}

func (s *Store) ChangeUserPassword(email string, password string, newPassword string) error {
	correctCredentials, err := s.AuthenticateUser(email, password)
	if err != nil {
		return err
	}
	if correctCredentials == false {
		return fmt.Errorf("Invalid user credentials for user: %s", email)
	}

	filter := bson.D{{Key: "email", Value: email}}

	update := bson.D{
		{Key: "$eq", Value: bson.D{
			{Key: "password", Value: newPassword},
		}},
	}

	userCollection := s.Collection("userDetails")

	updateResult, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return nil
}

func (s *Store) UpdateUserDetailsByEmail(email string, newDetails UserDetails) error {
	filter := bson.D{{Key: "email", Value: email}}

	update := bson.D{
		{Key: "$eq", Value: bson.D{
			{Key: "email", Value: newDetails.Email},
			{Key: "firstName", Value: newDetails.FirstName},
			{Key: "lastName", Value: newDetails.LastName},
			{Key: "gender", Value: newDetails.Gender},
			{Key: "dateOfBirth", Value: newDetails.DateOfBirth},
			{Key: "sendDeals", Value: newDetails.SendDeals},
		}},
	}

	userCollection := s.Collection("userDetails")

	updateResult, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return nil
}

func (s *Store) GetUserPassword(email string) (UserCredentials, error) {
	userCollection := s.Collection("userCredentials")

	var result UserCredentials

	err := userCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&result)
	if err != nil {
		return UserCredentials{}, err
	}

	return result, nil
}

func (s *Store) GetUserDetailsByEmail(email string) (UserDetails, error) {
	userCollection := s.Collection("userDetails")

	var result User

	err := userCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&result)
	if err != nil {
		return UserDetails{}, err
	}

	return userDetailsFromUser(result), nil
}

func (s *Store) GetAllUserDetails() ([]UserDetails, error) {
	findOptions := options.Find()
	userCollection := s.Collection("userDetails")

	var details []UserDetails

	cur, err := userCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return []UserDetails{}, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var elem User

		err := cur.Decode(&elem)
		if err != nil {
			return details, err
		}

		details = append(details, userDetailsFromUser(elem))
	}

	if err := cur.Err(); err != nil {
		return details, err
	}

	return details, nil
}
