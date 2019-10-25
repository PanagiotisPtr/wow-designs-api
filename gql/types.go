package gql

import "github.com/graphql-go/graphql"

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

var User = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"password": &graphql.Field{
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
