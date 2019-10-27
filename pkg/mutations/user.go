package mutations

import (
	"user-api/pkg/resolvers"

	"github.com/graphql-go/graphql"
)

func changePasswordMutation(resolver resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Change user password",
		Args: graphql.FieldConfigArgument{
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"newPassword": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: resolver.ChangePassword,
	}
}
