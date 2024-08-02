package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	schemainit "github.com/fnacarellidev/challenge-jbr/graphql-api/schema_init"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func FetchBackendCourtCase(p graphql.ResolveParams) (interface{}, error) {
	cnj, _ := p.Args["cnj"].(string)
	endpoint := os.Getenv("BACKEND_API_URL")+"fetch_court_case/"+cnj
	r, err := http.Get(endpoint)
	if err != nil {
		log.Println("GET at", endpoint, "failed with reason", err)
		return nil, nil
	}
	if r.StatusCode != http.StatusOK {
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
		"updates": courtCase.Updates,
	}, nil
}

func main() {
    h := handler.New(&handler.Config{
        Schema: schemainit.SchemaInit(FetchBackendCourtCase),
        Pretty: true,
    })

    http.Handle("/graphql", h)
    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

