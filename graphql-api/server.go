package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		bytes, _ := io.ReadAll(r.Body)
		return nil, errors.New(string(bytes))
	}

	var courtCase types.CourtCase
	err = json.NewDecoder(r.Body).Decode(&courtCase)
	if err != nil {
		return nil, err
	}

	courtOfOrigin, _ := p.Args["court_of_origin"].(string)
	if courtCase.CourtOfOrigin != courtOfOrigin {
		return nil, errors.New("no such case with cnj "+cnj+" at court "+courtOfOrigin)
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
	endpoint := os.Getenv("BACKEND_API_URL")+"register_court_case"
	b, _ := json.Marshal(courtCase)
	res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		bytes, _ := io.ReadAll(res.Body)
		return nil, errors.New(string(bytes))
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

func disableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
    h := handler.New(&handler.Config{
        Schema: schemainit.SchemaInit(FetchBackendCourtCase, AddCourtCase),
        Pretty: true,
    })

    http.Handle("/graphql", disableCors(h))
    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

