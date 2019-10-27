package gql

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql"
)

func ExecuteQuery(query string, schema graphql.Schema, userEmail string) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       context.WithValue(context.Background(), "userEmail", "userEmail"),
	})

	if len(result.Errors) > 0 {
		fmt.Printf("Unexpected errors inside ExecuteQuery: %v", result.Errors)
	}

	return result
}
