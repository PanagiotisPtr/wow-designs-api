package queries

import (
	"user-api/pkg/database"
	"user-api/pkg/resolvers"

	"github.com/graphql-go/graphql"
)

// NewQueryType returns a query type used to build
// the graphQL schema
func NewQueryType(s *database.Store) *graphql.Object {
	resolver := resolvers.Resolver{Store: s}

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user":         userQuery(resolver),
				"register":     registerQuery(resolver),
				"authenticate": authenticateQuery(resolver),
			},
		},
	)

	return queryType
}
