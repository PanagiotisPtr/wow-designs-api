package main

import (
	"context"
	"log"
	"net/http"
	"webserver-init/database"
	"webserver-init/gql"
	"webserver-init/server"

	"github.com/graphql-go/graphql"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func main() {
	router, store := initializeAPI()
	defer store.Close(context.TODO())

	log.Fatal(http.ListenAndServe(":4000", router))
}

func initializeAPI() (*chi.Mux, *database.Store) {
	router := chi.NewRouter()
	uri := "mongodb+srv://panagiotisptr:7ROLPv3AGdYCJwDk@cluster0-4pf3q.mongodb.net/test?retryWrites=true&w=majority"
	store, err := database.New(uri, "test")
	if err != nil {
		log.Fatal(err)
	}

	queryType := gql.NewQueryType(store)

	sc, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	s := server.Server{
		GqlSchema: &sc,
	}

	router.Use(
		render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.Logger,          // log api request calls
		middleware.DefaultCompress, // compress results, mostly gzipping assets and json
		middleware.StripSlashes,    // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,       // recover from panics without crashing server
	)

	router.Post("/api", s.GraphQL())

	return router, store
}
