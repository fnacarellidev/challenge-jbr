package tests

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/fnacarellidev/challenge-jbr/backend/endpoints"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/graphql-go/graphql"
	"github.com/julienschmidt/httprouter"
)

func FetchBackendCourtCase(p graphql.ResolveParams) (interface{}, error) {
	cnj, _ := p.Args["cnj"].(string)
	router := httprouter.New()
	router.GET("/fetch_court_case/:cnj", endpoints.FetchCourtCase)
	req, err := http.NewRequest("GET", "/fetch_court_case/"+cnj, nil)
	if err != nil {
		log.Println("ERROR", err)
	}

	var courtCase types.CourtCase
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	bytes, _ := io.ReadAll(rr.Body)
	err = json.Unmarshal(bytes, &courtCase)
	if err != nil {
		log.Println("failed unmarshal")
	}

	return map[string]interface{}{
		"cnj": courtCase.Cnj,
		"plaintiff": courtCase.Plaintiff,
		"defendant": courtCase.Defendant,
		"court_of_origin": courtCase.CourtOfOrigin,
		"start_date": courtCase.StartDate,
		"updates": courtCase.Updates,
	}, nil
}

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
		"updates": &graphql.Field{
			Type: graphql.NewList(caseUpdateType),
		},
    },
})

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
