package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fnacarellidev/challenge-jbr/backend/endpoints"
	"github.com/fnacarellidev/challenge-jbr/testutil"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Data struct {
    CourtCase types.CourtCase `json:"court_case"`
}

type GraphQLResponse struct {
    Data Data `json:"data"`
}

type GraphQLApiTestSuite struct {
	suite.Suite
	pgContainer *types.PostgresContainer
	ctx context.Context
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

func FetchBackendCourtCase(p graphql.ResolveParams) (interface{}, error) {
	cnj, _ := p.Args["cnj"].(string)
	router := httprouter.New()
	router.GET("/fetch_court_case/:cnj", endpoints.FetchCourtCase)
	endpoint := os.Getenv("BACKEND_API_URL")+"fetch_court_case/"+cnj
	req, err := http.NewRequest("GET", endpoint, nil)
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

func (suite *GraphQLApiTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testutil.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer
	os.Setenv("DATABASE_URL", suite.pgContainer.ConnectionString)
	os.Setenv("BACKEND_API_URL", "/")
}

func (suite *GraphQLApiTestSuite) TearDownTestSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *GraphQLApiTestSuite) TestFetchCourtCaseAliceBob() {
	t := suite.T()
	router := httprouter.New()
	router.GET("/fetch_court_case/:cnj", endpoints.FetchCourtCase)

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})
	router.Handler("POST", "/graphql", h)
	query := ` 
	{
		"query": "query($cnj: String!) { court_case(cnj: $cnj) { cnj plaintiff defendant court_of_origin start_date updates { update_date update_details } } }",
		"variables": {
			"cnj": "5001682-88.2024.8.13.0672"
		}
	}
	`
	jsonStr := []byte(query)
	req, err := http.NewRequest("POST", "/graphql", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Errorf("GET request failed: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	bytes, err := io.ReadAll(rr.Body)
	if err != nil {
		log.Println("failed to read bytes", err)
	}

	var graphQLResponse GraphQLResponse
	err = json.Unmarshal(bytes, &graphQLResponse)
	if err != nil {
		log.Println("failed to unmarshal graphql", err)
	}

	courtCase := graphQLResponse.Data.CourtCase
	assert.Equal(t, "Alice Johnson", courtCase.Plaintiff, "Plaintiff name does not match")
}

func TestGraphQLApiSuite(t *testing.T) {
    suite.Run(t, new(GraphQLApiTestSuite))
}
