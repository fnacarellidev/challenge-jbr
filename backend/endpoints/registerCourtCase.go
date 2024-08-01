package endpoints

import (
	"os"
	"log"
	"context"
	"encoding/json"
	"net/http"

	"github.com/fnacarellidev/challenge-jbr/backend/.sqlcbuild/pgquery"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
)

func RegisterCourtCase(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to db")
	}
	defer conn.Close(context.Background())

	var courtCase types.CourtCase
	err = json.NewDecoder(r.Body).Decode(&courtCase)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	sqlc := pgquery.New(conn)
	_, err = sqlc.InsertCourtCase(context.Background(), pgquery.InsertCourtCaseParams{
		Cnj: courtCase.Cnj,
		Plaintiff: courtCase.Plaintiff,
		Defendant: courtCase.Defendant,
		CourtOfOrigin: courtCase.CourtOfOrigin,
		StartDate: pgtype.Timestamptz{
			Time: courtCase.StartDate,
			Valid: true,
		},
	})
	if err != nil {
		http.Error(w, "Failed to save on db", http.StatusBadRequest)
		return
	}
}

