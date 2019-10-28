package mutations

import (
	"user-api/pkg/database"
	"user-api/pkg/resolvers"

	"github.com/graphql-go/graphql"
)

// NewMutationType returns a mutation type used to
// build the graphQL schema
func NewMutationType(s *database.Store) *graphql.Object {
	resolver := resolvers.Resolver{Store: s}

	var mutationType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"register":          registerMutation(resolver),
				"changePassword":    changePasswordMutation(resolver),
				"changeUserDetails": changeUserDetails(resolver),
				"terminate":         terminateMutation(resolver),
			},
		},
	)

	return mutationType
}
