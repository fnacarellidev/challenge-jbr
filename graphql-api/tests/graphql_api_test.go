package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/fnacarellidev/challenge-jbr/backend/endpoints"
	schemainit "github.com/fnacarellidev/challenge-jbr/graphql-api/schema_init"
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
	router *httprouter.Router
}

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

	if rr.Result().StatusCode != http.StatusOK {
		return nil, nil
	}

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

func AddCourtCase(p graphql.ResolveParams) (interface{}, error) {
	router := httprouter.New()
	router.POST("/register_court_case", endpoints.RegisterCourtCase)
	courtCase := types.CourtCase{
		Cnj: p.Args["cnj"].(string),
		Plaintiff: p.Args["plaintiff"].(string),
		Defendant: p.Args["defendant"].(string),
		CourtOfOrigin: p.Args["court_of_origin"].(string),
		StartDate: p.Args["start_date"].(time.Time),
	}

	genericUpdates := p.Args["updates"].([]interface{})
	for _, genericUpdate := range genericUpdates {
		updateMap := genericUpdate.(map[string]interface{})
		updateDate := updateMap["update_date"].(time.Time)
		updateDetails := updateMap["update_details"].(string)
		courtCase.Updates = append(courtCase.Updates, types.CaseUpdate{
			UpdateDate:    updateDate,
			UpdateDetails: updateDetails,
		})
	}
	endpoint := "/register_court_case"
	b, _ := json.Marshal(courtCase)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		bytes, _ := io.ReadAll(rr.Result().Body)
		return nil, errors.New(string(bytes))
	}

	return map[string]interface{}{
		"cnj": courtCase.Cnj,
	}, nil
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
		Schema: schemainit.SchemaInit(FetchBackendCourtCase, AddCourtCase),
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

func (suite *GraphQLApiTestSuite) TestFetchCourtCaseAliceBobPlaintiffOnly() {
	t := suite.T()
	expectedPlaintiff := "Alice Johnson"
	query := ` 
	{
		"query": "query($cnj: String!) { court_case(cnj: $cnj) { plaintiff } }",
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

	var graphQLResponse struct {
		Data struct {
			CourtCase struct {
				Plaintiff string `json:"plaintiff"`
			} `json:"court_case"`
		} `json:"data"`
	}
	err = json.Unmarshal(bytes, &graphQLResponse)
	if err != nil {
		log.Println("failed to unmarshal graphql", err)
	}

	courtCase := graphQLResponse.Data.CourtCase
	assert.Equal(t, expectedPlaintiff, courtCase.Plaintiff, "plaintiff does not match")
}

func (suite *GraphQLApiTestSuite) TestFetchCourtCaseAliceBobPlaintiffDefendant() {
	t := suite.T()
	expectedPlaintiff := "Alice Johnson"
	expectedDefendant := "Bob Smith"
	query := ` 
	{
		"query": "query($cnj: String!) { court_case(cnj: $cnj) { plaintiff defendant } }",
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

	var graphQLResponse struct {
		Data struct {
			CourtCase struct {
				Plaintiff string `json:"plaintiff"`
				Defendant string `json:"defendant"`
			} `json:"court_case"`
		} `json:"data"`
	}
	err = json.Unmarshal(bytes, &graphQLResponse)
	if err != nil {
		log.Println("failed to unmarshal graphql", err)
	}

	courtCase := graphQLResponse.Data.CourtCase
	assert.Equal(t, expectedPlaintiff, courtCase.Plaintiff, "plaintiff does not match")
	assert.Equal(t, expectedDefendant, courtCase.Defendant, "defendant does not match")
}

func (suite *GraphQLApiTestSuite) TestFetchCourtThatDoesntExist() {
	t := suite.T()
	query := `
	{
		"query": "query($cnj: String!) { court_case(cnj: $cnj) { plaintiff defendant } }",
		"variables": {
			"cnj": "courtcasethatdoesntexist"
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

	assert.Equal(t, graphQLResponse.Data.CourtCase, types.CourtCase{}, "response object was not 'null'")
}

func (suite *GraphQLApiTestSuite) TestFetchOnlyUpdatesFromCourtCase() {
	t := suite.T()
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
		"query": "query($cnj: String!) { court_case(cnj: $cnj) { updates { update_date update_details } } }",
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

	var graphQLResponse struct {
		Data struct {
			CourtCase struct {
				Updates []types.CaseUpdate `json:"updates"`
			} `json:"court_case"`
		} `json:"data"`
	}
	err = json.Unmarshal(bytes, &graphQLResponse)
	if err != nil {
		log.Println("failed to unmarshal graphql", err)
	}

	courtCase := graphQLResponse.Data.CourtCase

	assert.Equal(t, expectedUpdates[1], courtCase.Updates[1].UpdateDetails, "Second update does not match")
	assert.Equal(t, expectedUpdates[2], courtCase.Updates[2].UpdateDetails, "Third update does not match")
	assert.Equal(t, expectedUpdates[0], courtCase.Updates[0].UpdateDetails, "first update does not match")
	assert.Equal(t, expectedUpdatesDates[0], courtCase.Updates[0].UpdateDate, "First update date does not match")
	assert.Equal(t, expectedUpdatesDates[1], courtCase.Updates[1].UpdateDate, "Second update date does not match")
	assert.Equal(t, expectedUpdatesDates[2], courtCase.Updates[2].UpdateDate, "Third update date does not match")
}

func (suite *GraphQLApiTestSuite) TestInsertCourtCase() {
	t := suite.T()
	query := `
	{
		"query": "mutation new_court_case($cnj: String!, $plaintiff: String!, $defendant: String!, $court_of_origin: String!, $start_date: DateTime!, $updates: [CaseUpdateInput]) { new_court_case(cnj: $cnj, plaintiff: $plaintiff, defendant: $defendant, court_of_origin: $court_of_origin, start_date: $start_date, updates: $updates) { cnj } }",
		"variables": {
			"cnj": "12345-67.2024.8.1.0001",
			"plaintiff": "John Doe",
			"defendant": "Foo Bar",
			"court_of_origin": "First Court",
			"start_date": "2024-01-01T00:00:00Z",
			"updates": [
			{
				"update_date": "2024-02-01T00:00:00Z",
				"update_details": "Initial hearing scheduled"
			},
			{
				"update_date": "2024-03-01T00:00:00Z",
				"update_details": "Preliminary ruling issued"
			}
			]
		}
	}
	`
	jsonStr := []byte(query)
	req, _ := http.NewRequest("POST", "/graphql", bytes.NewBuffer(jsonStr))
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)
	bytes, err := io.ReadAll(rr.Body)
	if err != nil {
		log.Println("failed to read bytes", err)
	}

	var graphQLResponse struct {
		Data struct {
			NewCourtCase struct {
				Cnj string `json:"cnj"`
			} `json:"new_court_case"`
		} `json:"data"`
	}
	err = json.Unmarshal(bytes, &graphQLResponse)
	if err != nil {
		log.Println("failed to unmarshal graphql", err)
	}

	newCourtCase := graphQLResponse.Data.NewCourtCase
	assert.Equal(t, "12345-67.2024.8.1.0001", newCourtCase.Cnj, "new cnj did not match expected")
}

func (suite *GraphQLApiTestSuite) TestInsertCourtCaseThatAlreadyExist() {
	t := suite.T()
	query := `
	{
		"query": "mutation new_court_case($cnj: String!, $plaintiff: String!, $defendant: String!, $court_of_origin: String!, $start_date: DateTime!, $updates: [CaseUpdateInput]) { new_court_case(cnj: $cnj, plaintiff: $plaintiff, defendant: $defendant, court_of_origin: $court_of_origin, start_date: $start_date, updates: $updates) { cnj } }",
		"variables": {
			"cnj": "5001682-88.2024.8.13.0672",
			"plaintiff": "John Doe",
			"defendant": "Foo Bar",
			"court_of_origin": "First Court",
			"start_date": "2024-01-01T00:00:00Z",
			"updates": [
			{
				"update_date": "2024-02-01T00:00:00Z",
				"update_details": "Initial hearing scheduled"
			},
			{
				"update_date": "2024-03-01T00:00:00Z",
				"update_details": "Preliminary ruling issued"
			}
			]
		}
	}
	`
	jsonStr := []byte(query)
	req, _ := http.NewRequest("POST", "/graphql", bytes.NewBuffer(jsonStr))
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)
	bytes, err := io.ReadAll(rr.Body)
	if err != nil {
		log.Println("failed to read bytes", err)
	}

	var graphQLResponse struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	err = json.Unmarshal(bytes, &graphQLResponse)
	if err != nil {
		log.Println("failed to unmarshal graphql", err)
	}

	errors := graphQLResponse.Errors

	assert.GreaterOrEqual(t, 1, len(errors), "amount of errors was not >= 1")
	assert.Equal(t, "case already exists", graphQLResponse.Errors[0].Message, "error message did not match")
}

func TestGraphQLApiSuite(t *testing.T) {
    suite.Run(t, new(GraphQLApiTestSuite))
}
