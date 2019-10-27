package mutations

import (
	"user-api/pkg/resolvers"
	"user-api/pkg/types"

	"github.com/graphql-go/graphql"
)

func changePasswordMutation(resolver resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        types.UserDetails,
		Description: "Create new user",
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	}
}
