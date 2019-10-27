package mutations

import (
	"user-api/pkg/resolvers"

	"github.com/graphql-go/graphql"
)

func changePasswordMutation(resolver resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Create new user",
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: resolver.ChangePassword,
	}
}
