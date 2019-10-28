package queries

import (
	"user-api/pkg/resolvers"
	"user-api/pkg/types"

	"github.com/graphql-go/graphql"
)

func userDetailsQuery(resolver resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        types.UserDetails,
		Description: "Get user details using JWT in cookie",
		Resolve:     resolver.UserDetails,
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
		Resolve: resolver.Authenticate,
	}
}
