package queries

import (
	"user-api/pkg/database"
	"user-api/pkg/resolvers"
	"user-api/pkg/types"

	"github.com/graphql-go/graphql"
)

func NewQueryType(s *database.Store) *graphql.Object {
	resolver := resolvers.Resolver{Store: s}

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type:        types.UserDetails,
					Description: "Get user details using JWT in cookie",
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
					Type:        types.SessionToken,
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
			},
		},
	)

	return queryType
}
