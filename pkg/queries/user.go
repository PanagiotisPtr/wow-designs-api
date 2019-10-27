package queries

import (
	"user-api/pkg/resolvers"
	"user-api/pkg/types"

	"github.com/graphql-go/graphql"
)

func userQuery(resolver resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        types.UserDetails,
		Description: "Get user details using JWT in cookie",
		Resolve:     resolver.UserDetailsResolver,
	}
}

func signupQuery(resolver resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
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
	}
}

func authenticateQuery(resolver resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        types.SessionToken,
		Description: "Authenticate user and get token",
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: resolver.AuthenticationResolver,
	}
}
