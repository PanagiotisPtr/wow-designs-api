package mutations

import (
	"user-api/pkg/resolvers"

	"github.com/graphql-go/graphql"
)

func registerMutation(resolver resolvers.Resolver) *graphql.Field {
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
		Resolve: resolver.Register,
	}
}

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

func changeUserDetails(resolver resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Change user details",
		Args: graphql.FieldConfigArgument{
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
		Resolve: resolver.ChangeUserDetails,
	}
}

func terminateMutation(resolver resolvers.Resolver) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Terminate current user account",
		Args: graphql.FieldConfigArgument{
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: resolver.Terminate,
	}
}
