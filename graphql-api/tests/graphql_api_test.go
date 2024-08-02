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
	"time"

	"github.com/fnacarellidev/challenge-jbr/testutil"
	"github.com/fnacarellidev/challenge-jbr/types"
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
	router *httprouter.Router
}

func (suite *GraphQLApiTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testutil.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer
	suite.router = httprouter.New()
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})
	suite.router.Handler("POST", "/graphql", h)
	os.Setenv("DATABASE_URL", suite.pgContainer.ConnectionString)
}

func (suite *GraphQLApiTestSuite) TearDownTestSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *GraphQLApiTestSuite) TestFetchCourtCaseAliceBobAllInfo() {
	t := suite.T()
	expectedCnj := "5001682-88.2024.8.13.0672"
	expectedPlaintiff := "Alice Johnson"
	expectedDefendant := "Bob Smith"
	courtOfOrigin := "TJSP"
	expectedUpdates := []string{
		"Defendant requested a delay for response.",
		"Plaintiff submitted additional evidence.",
		"Initial hearing scheduled for August 15, 2024.",
	}
	expectedUpdatesDates := []time.Time{
		time.Date(2024, 8, 2, 6, 0, 0, 0, time.Local),
		time.Date(2024, 8, 1, 11, 30, 0, 0, time.Local),
		time.Date(2024, 7, 31, 7, 0, 0, 0, time.Local),
	}
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
	suite.router.ServeHTTP(rr, req)
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
	assert.Equal(t, expectedCnj, courtCase.Cnj, "cnj does not match")
	assert.Equal(t, expectedPlaintiff, courtCase.Plaintiff, "plaintiff does not match")
	assert.Equal(t, expectedDefendant, courtCase.Defendant, "defendant does not match")
	assert.Equal(t, courtOfOrigin, courtCase.CourtOfOrigin, "court of origin does not match")
	assert.Equal(t, expectedUpdates[1], courtCase.Updates[1].UpdateDetails, "Second update does not match")
	assert.Equal(t, expectedUpdates[2], courtCase.Updates[2].UpdateDetails, "Third update does not match")
	assert.Equal(t, expectedUpdates[0], courtCase.Updates[0].UpdateDetails, "first update does not match")
	assert.Equal(t, expectedUpdatesDates[0], courtCase.Updates[0].UpdateDate, "First update date does not match")
	assert.Equal(t, expectedUpdatesDates[1], courtCase.Updates[1].UpdateDate, "Second update date does not match")
	assert.Equal(t, expectedUpdatesDates[2], courtCase.Updates[2].UpdateDate, "Third update date does not match")
}

func TestGraphQLApiSuite(t *testing.T) {
    suite.Run(t, new(GraphQLApiTestSuite))
}
