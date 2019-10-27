package gql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

func ExecuteQuery(query string, schema graphql.Schema, cookie *http.Cookie) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       context.WithValue(context.Background(), "cookie", cookie),
	})

	if len(result.Errors) > 0 {
		fmt.Printf("Unexpected errors inside ExecuteQuery: %v", result.Errors)
	}

	return result
}
