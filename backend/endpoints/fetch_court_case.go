package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/fnacarellidev/challenge-jbr/backend/.sqlcbuild/pgquery"
	"github.com/fnacarellidev/challenge-jbr/backend/endpoints/utils"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

func FetchUpdatesFromCase(cnj string, sqlc *pgquery.Queries) ([]types.CaseUpdate, error) {
	rows, err := sqlc.GetCaseUpdates(context.Background(), cnj)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	var caseUpdates []types.CaseUpdate
	for _, row := range(rows) {
		caseUpdate := types.CaseUpdate{
			UpdateDate: row.UpdateDate.Time,
			UpdateDetails: row.UpdateDetails,
		}
		caseUpdates = append(caseUpdates, caseUpdate)
	}
	return caseUpdates, nil
}

// todo think about case where the court has already been registered
func FetchCourtCase(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to db")
	}
	defer conn.Close(context.Background())

	cnjLookup := ps.ByName("cnj")
	sqlc := pgquery.New(conn)
	courtCase, error := sqlc.GetCourtCase(context.Background(), cnjLookup)
	if error != nil {
		if error == pgx.ErrNoRows {
			utils.SendError(w, "no case with cnj "+cnjLookup, http.StatusNotFound)
		} else {
			utils.SendError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	caseUpdates, err := FetchUpdatesFromCase(cnjLookup, sqlc)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBytes, _ := json.Marshal(types.CourtCase{
		Cnj: courtCase.Cnj,
		Plaintiff: courtCase.Plaintiff,
		Defendant: courtCase.Defendant,
		CourtOfOrigin: courtCase.CourtOfOrigin,
		StartDate: courtCase.StartDate.Time,
		Updates: caseUpdates,
	})
	w.Write(jsonBytes)
}

