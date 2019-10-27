package types

import "github.com/graphql-go/graphql"

// UserDetails object contains information about a user.
// This information does not include their password
var UserDetails = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
			"gender": &graphql.Field{
				Type: graphql.String,
			},
			"dateOfBirth": &graphql.Field{
				Type: graphql.String,
			},
			"sendDeals": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)
