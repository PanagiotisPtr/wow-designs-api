package gql

import (
	"webserver-init/database"

	"github.com/graphql-go/graphql"
)

func NewQueryType(s *database.Store) *graphql.Object {
	resolver := Resolver{store: s}

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type:        UserDetails,
					Description: "Get user by email",
					Resolve:     resolver.UserDetailsResolver,
				},
				"signup": &graphql.Field{
					Type:        graphql.Boolean,
					Description: "Register new user",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"firstName": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"lastName": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"gender": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"dateOfBirth": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"sendDeals": &graphql.ArgumentConfig{
							Type: graphql.Boolean,
						},
					},
					Resolve: resolver.RegisterResolver,
				},
				"authenticate": &graphql.Field{
					Type:        SessionToken,
					Description: "Authenticate user and get token",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolver.AuthenticationResolver,
				},
				"list": &graphql.Field{
					Type:        graphql.NewList(UserDetails),
					Description: "Get list of all users",
					Resolve:     resolver.ListResolver,
				},
			},
		},
	)

	return queryType
}
