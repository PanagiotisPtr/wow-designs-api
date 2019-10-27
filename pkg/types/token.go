package types

import "github.com/graphql-go/graphql"

// SessionToken contains a field token which is the JSON Web Token
// used for identifying a user
var SessionToken = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "SessionToken",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
