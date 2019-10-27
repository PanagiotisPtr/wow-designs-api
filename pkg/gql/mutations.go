package gql

import (
	"user-api/pkg/database"

	"github.com/graphql-go/graphql"
)

func NewMutationType(s *database.Store) *graphql.Object {
	//resolver := Resolver{store: s}

	var mutationType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"fields": &graphql.Field{
					Type:        User,
					Description: "Create new user",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"password": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						}, // todo
					},
				},
			},
		},
	)

	return mutationType
}
