package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/graphql-go/graphql"
)

var caseUpdateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CaseUpdate",
	Fields: graphql.Fields{
		"update_date": &graphql.Field{
			Type: graphql.DateTime,
		},
		"update_details": &graphql.Field{
			Type: graphql.String,
		},
	},
})

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
	endpoint := os.Getenv("BACKEND_API_URL")+"fetch_court_case/"+cnj
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

func FetchBackendCaseUpdate(p graphql.ResolveParams) (interface{}, error) {
	cnj, _ := p.Args["cnj"].(string)
	endpoint := os.Getenv("BACKEND_API_URL")+"fetch_updates_from_case/"+cnj
	r, err := http.Get(endpoint)
	if err != nil {
		log.Println("GET at", endpoint, "failed with reason", err)
		return nil, nil
	}

	var caseUpdates []types.CaseUpdate
	err = json.NewDecoder(r.Body).Decode(&caseUpdates)
	if err != nil {
		log.Println("error:", err)
		return nil, nil
	}

	return caseUpdates, nil
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
		"case_updates": &graphql.Field{
			Type: graphql.NewList(caseUpdateType),
			Args: graphql.FieldConfigArgument{
				"cnj": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: FetchBackendCaseUpdate,
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: rootQuery,
})
