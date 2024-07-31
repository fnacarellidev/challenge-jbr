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
        "court_of_origin": &graphql.Field{
            Type: graphql.String,
        },
        "start_date": &graphql.Field{
            Type: graphql.DateTime,
        },
    },
})

func FetchBackendCourtCase(p graphql.ResolveParams) (interface{}, error) {
	cnj, _ := p.Args["cnj"].(string)
	endpoint := "http://localhost:8081/fetch_court_case/"+cnj
	r, err := http.Get(endpoint)
	if err != nil {
		log.Println("GET at", endpoint, "failed with reason", err)
		return nil, nil
	}

	var courtCase types.CourtCase
	err = json.NewDecoder(r.Body).Decode(&courtCase)
	if err != nil {
		return nil, nil
	}

	return map[string]interface{}{
		"cnj": courtCase.Cnj,
		"plaintiff": courtCase.Plaintiff,
		"defendant": courtCase.Defendant,
		"court_of_origin": courtCase.CourtOfOrigin,
		"start_date": courtCase.StartDate,
	}, nil
}

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"court_case": &graphql.Field{
			Type: courtCaseType,
			Args: graphql.FieldConfigArgument{
				"cnj": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: FetchBackendCourtCase,
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: rootQuery,
})
