package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/fnacarellidev/challenge-jbr/backend/endpoints"
	"github.com/fnacarellidev/challenge-jbr/testutil"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BackendApiTestSuite struct {
	suite.Suite
	pgContainer *types.PostgresContainer
	ctx context.Context
}

func (suite *BackendApiTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testutil.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer
	os.Setenv("DATABASE_URL", suite.pgContainer.ConnectionString)
}

func (suite *BackendApiTestSuite) TearDownTestSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *BackendApiTestSuite) TestFetchCourtCaseAliceBob() {
	t := suite.T()
	router := httprouter.New()
	expectedPlaintiff := "Alice Johnson"
	expectedDefendant := "Bob Smith"
	expectedCourtOfOrigin := "TJSP"
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
	endpoint := "/fetch_court_case/5001682-88.2024.8.13.0672"
	router.GET("/fetch_court_case/:cnj", endpoints.FetchCourtCase)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Errorf("GET request failed: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code of %v, got %v", http.StatusOK, rr.Code)
    }

	var courtCase types.CourtCase
	err = json.Unmarshal(rr.Body.Bytes(), &courtCase)
	if err != nil {
		t.Errorf("Failed to unmarshal api response: %v", err)
	}

	assert.Equal(t, expectedPlaintiff, courtCase.Plaintiff, "Plaintiff name does not match")
	assert.Equal(t, expectedDefendant, courtCase.Defendant, "Defendant name does not match")
	assert.Equal(t, expectedCourtOfOrigin, courtCase.CourtOfOrigin, "Court of Origin does not match")
	assert.Equal(t, expectedUpdates[0], courtCase.Updates[0].UpdateDetails, "First update does not match")
	assert.Equal(t, expectedUpdates[1], courtCase.Updates[1].UpdateDetails, "Second update does not match")
	assert.Equal(t, expectedUpdates[2], courtCase.Updates[2].UpdateDetails, "Third update does not match")
	assert.Equal(t, expectedUpdatesDates[0], courtCase.Updates[0].UpdateDate, "First update date does not match")
	assert.Equal(t, expectedUpdatesDates[1], courtCase.Updates[1].UpdateDate, "Second update date does not match")
	assert.Equal(t, expectedUpdatesDates[2], courtCase.Updates[2].UpdateDate, "Third update date does not match")
}

func (suite *BackendApiTestSuite) TestFetchCourtCaseMichaelSarah() {
	t := suite.T()
	router := httprouter.New()
	expectedPlaintiff := "Michael Brown"
	expectedDefendant := "Sarah Davis"
	expectedCourtOfOrigin := "TJSP"
	expectedUpdates := []string{
		"Hearing date scheduled for August 10, 2024.",
		"Defendant’s lawyer filed a motion for dismissal.",
		"Witness statements collected.",
		"Case file reviewed by judge.",
	}
	expectedUpdatesDates := []time.Time{
		time.Date(2024, 8, 3, 7, 0, 0, 0, time.Local),
		time.Date(2024, 8, 2, 10, 45, 0, 0, time.Local),
		time.Date(2024, 8, 1, 6, 30, 0, 0, time.Local),
		time.Date(2024, 7, 30, 8, 0, 0, 0, time.Local),
	}
	endpoint := "/fetch_court_case/3562061-02.2024.8.13.0431"
	router.GET("/fetch_court_case/:cnj", endpoints.FetchCourtCase)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Errorf("GET request failed: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code of %v, got %v", http.StatusOK, rr.Code)
    }

	var courtCase types.CourtCase
	err = json.Unmarshal(rr.Body.Bytes(), &courtCase)
	if err != nil {
		t.Errorf("Failed to unmarshal api response: %v", err)
	}

	assert.Equal(t, expectedPlaintiff, courtCase.Plaintiff, "Plaintiff name does not match")
	assert.Equal(t, expectedDefendant, courtCase.Defendant, "Defendant name does not match")
	assert.Equal(t, expectedCourtOfOrigin, courtCase.CourtOfOrigin, "Court of Origin does not match")
	assert.Equal(t, expectedUpdates[0], courtCase.Updates[0].UpdateDetails, "First update does not match")
	assert.Equal(t, expectedUpdates[1], courtCase.Updates[1].UpdateDetails, "Second update does not match")
	assert.Equal(t, expectedUpdates[2], courtCase.Updates[2].UpdateDetails, "Third update does not match")
	assert.Equal(t, expectedUpdates[3], courtCase.Updates[3].UpdateDetails, "Fourth update does not match")
	assert.Equal(t, expectedUpdatesDates[0], courtCase.Updates[0].UpdateDate, "First update date does not match")
	assert.Equal(t, expectedUpdatesDates[1], courtCase.Updates[1].UpdateDate, "Second update date does not match")
	assert.Equal(t, expectedUpdatesDates[2], courtCase.Updates[2].UpdateDate, "Third update date does not match")
	assert.Equal(t, expectedUpdatesDates[3], courtCase.Updates[3].UpdateDate, "Fourth update date does not match")
}

func (suite *BackendApiTestSuite) TestFetchCourtCaseThatDoesntExist() {
	t := suite.T()
	router := httprouter.New()
	expectedResponse := "no case with cnj casethatdoesntexist"
	endpoint := "/fetch_court_case/casethatdoesntexist"
	router.GET("/fetch_court_case/:cnj", endpoints.FetchCourtCase)

	req, _ := http.NewRequest("GET", endpoint, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	var errResponse types.ErrResponse
	json.Unmarshal(rr.Body.Bytes(), &errResponse)

	assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode, "Status was not 404.")
	assert.Equal(t, expectedResponse, errResponse.Error, "Response body was not '404 not found'")
}

func (suite *BackendApiTestSuite) TestInsertCourtCaseThatExists() {
	t := suite.T()
	router := httprouter.New()
	expectedResponse := "case already exists"
	endpoint := "/register_court_case"
	router.POST(endpoint, endpoints.RegisterCourtCase)

	jsonBody, _ := json.Marshal(
		types.CourtCase{
			Cnj: "5001682-88.2024.8.13.0672",
			Plaintiff: "John Doe",
			Defendant: "Foo Bar",
			StartDate: time.Now().Local(),
			CourtOfOrigin: "FOO",
		},
	)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	var errResponse types.ErrResponse
	json.Unmarshal(rr.Body.Bytes(), &errResponse)

	assert.Equal(t, http.StatusConflict, rr.Result().StatusCode, "Status does not match.")
	assert.Equal(t, expectedResponse, errResponse.Error, "Err message does not match.")
}

func (suite *BackendApiTestSuite) TestInsertInvalidCourtCase() {
	t := suite.T()
	router := httprouter.New()
	expectedResponse := "invalid request payload"
	endpoint := "/register_court_case"
	router.POST(endpoint, endpoints.RegisterCourtCase)

	jsonBody, _ := json.Marshal(
		struct {
			Cnj int `json:"cnj"`
		}{
			Cnj: 1,
		},
	)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	var errResponse types.ErrResponse
	json.Unmarshal(rr.Body.Bytes(), &errResponse)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode, "Status does not match.")
	assert.Equal(t, expectedResponse, errResponse.Error, "Err message does not match.")
}

func (suite *BackendApiTestSuite) TestInsertCourtCaseEmptyCnj() {
	t := suite.T()
	router := httprouter.New()
	expectedResponse := "cnj field cannot be empty"
	endpoint := "/register_court_case"
	router.POST(endpoint, endpoints.RegisterCourtCase)

	jsonBody, _ := json.Marshal(
		types.CourtCase{
			Cnj: "",
			Plaintiff: "John Doe",
			Defendant: "Foo Bar",
			StartDate: time.Now().UTC(),
			CourtOfOrigin: "FOO",
		},
	)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	var errResponse types.ErrResponse
	json.Unmarshal(rr.Body.Bytes(), &errResponse)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode, "Status does not match.")
	assert.Equal(t, expectedResponse, errResponse.Error, "Err message does not match.")
}

func TestBackendApiSuite(t *testing.T) {
    suite.Run(t, new(BackendApiTestSuite))
}
