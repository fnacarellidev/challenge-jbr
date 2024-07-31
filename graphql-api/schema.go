package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/fnacarellidev/challenge-jbr/types"
)

var courtCaseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CourtCase",
    Fields: graphql.Fields{
        "cnj": &graphql.Field{
            Type: graphql.String,
        },
        "plaintiff": &graphql.Field{
            Type: graphql.String,
        },
        "defendant": &graphql.Field{
            Type: graphql.String,
        },
        "courtOfOrigin": &graphql.Field{
            Type: graphql.String,
        },
        "startDate": &graphql.Field{
            Type: graphql.DateTime,
        },
    },
})

func FetchBackendCourtCase(p graphql.ResolveParams) (interface{}, error) {
	r, err := http.Post("http://localhost:8081/fetch_court_case/123", "application/json", nil)
	if err != nil {
		log.Println("failed", err)
		return nil, nil
	}

	var courtCase types.CourtCase
	err = json.NewDecoder(r.Body).Decode(&courtCase)
	if err != nil {
		log.Println("failed 2", err)
		return nil, nil
	}

	return map[string]interface{}{
		"cnj": courtCase.Cnj,
		"plaintiff": courtCase.Plaintiff,
		"defendant": courtCase.Defendant,
		"courtOfOrigin": courtCase.CourtOfOrigin,
		"startDate": courtCase.StartDate,
	}, nil
}

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"courtCase": &graphql.Field{
			Type: courtCaseType,
			Resolve: FetchBackendCourtCase,
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: rootQuery,
})
