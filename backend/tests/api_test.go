package tests

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/fnacarellidev/challenge-jbr/backend/endpoints"
	"github.com/fnacarellidev/challenge-jbr/backend/tests/testhelpers"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BackendApiTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	ctx context.Context
}

func (suite *BackendApiTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
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
}

func (suite *BackendApiTestSuite) TestFetchCourtCaseMichaelSarah() {
	t := suite.T()
	router := httprouter.New()
	expectedPlaintiff := "Michael Brown"
	expectedDefendant := "Sarah Davis"
	expectedCourtOfOrigin := "TJSP"
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
}

func (suite *BackendApiTestSuite) TestFetchUpdatesOnAliceVersusBob() {
	t := suite.T()
	router := httprouter.New()
	expectedUpdates := []string{
		"Defendant requested a delay for response.",
		"Plaintiff submitted additional evidence.",
		"Initial hearing scheduled for August 15, 2024.",
	}
	expectedDates := []time.Time{
		time.Date(2024, 8, 2, 6, 0, 0, 0, time.Local),
		time.Date(2024, 8, 1, 11, 30, 0, 0, time.Local),
		time.Date(2024, 7, 31, 7, 0, 0, 0, time.Local),
	}
	endpoint := "/fetch_updates_from_case/5001682-88.2024.8.13.0672"
	router.GET("/fetch_updates_from_case/:cnj", endpoints.FetchUpdatesFromCase)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Errorf("GET request failed: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code of %v, got %v", http.StatusOK, rr.Code)
    }

	var caseUpdates []types.CaseUpdate
	err = json.Unmarshal(rr.Body.Bytes(), &caseUpdates)
	if err != nil {
		t.Errorf("Failed to unmarshal api response: %v", err)
	}

	for i, update := range caseUpdates {
		assert.Equal(t, expectedUpdates[i], update.UpdateDetails, "Case update detail does not match")
		assert.Equal(t, expectedDates[i], update.UpdateDate, "Case update date does not match")
	}
}

func (suite *BackendApiTestSuite) TestFetchUpdatesOnChrisVersusJessica() {
	t := suite.T()
	router := httprouter.New()
	expectedUpdates := []string{
		"Hearing date scheduled for August 10, 2024.",
		"Defendant’s lawyer filed a motion for dismissal.",
		"Witness statements collected.",
		"Case file reviewed by judge.",
	}
	expectedDates := []time.Time{
		time.Date(2024, 8, 3, 7, 0, 0, 0, time.Local),
		time.Date(2024, 8, 2, 10, 45, 0, 0, time.Local),
		time.Date(2024, 8, 1, 6, 30, 0, 0, time.Local),
		time.Date(2024, 7, 30, 8, 0, 0, 0, time.Local),
	}
	endpoint := "/fetch_updates_from_case/6772130-04.2024.8.13.0161"
	router.GET("/fetch_updates_from_case/:cnj", endpoints.FetchUpdatesFromCase)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Errorf("GET request failed: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code of %v, got %v", http.StatusOK, rr.Code)
    }

	var caseUpdates []types.CaseUpdate
	err = json.Unmarshal(rr.Body.Bytes(), &caseUpdates)
	if err != nil {
		t.Errorf("Failed to unmarshal api response: %v", err)
	}

	for i, update := range caseUpdates {
		assert.Equal(t, expectedUpdates[i], update.UpdateDetails, "Case update detail does not match")
		assert.Equal(t, expectedDates[i], update.UpdateDate, "Case update date does not match")
	}
}

func TestBackendApiSuite(t *testing.T) {
    suite.Run(t, new(BackendApiTestSuite))
}
