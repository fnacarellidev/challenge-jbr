package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fnacarellidev/challenge-jbr/backend/endpoints"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BackendApiTestSuite struct {
	suite.Suite
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


func TestBackendApiSuite(t *testing.T) {
    suite.Run(t, new(BackendApiTestSuite))
}
