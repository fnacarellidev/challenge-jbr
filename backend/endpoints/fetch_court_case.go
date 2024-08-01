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
	"github.com/julienschmidt/httprouter"
)

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
		var errorResponse types.ErrResponse
		errorResponse.Error = "internal server error"
		statusCode := http.StatusInternalServerError

		if error == pgx.ErrNoRows {
			statusCode = http.StatusNotFound
			errorResponse.Error = "no case with cnj "+cnjLookup
		}
		bytes, _ := json.Marshal(errorResponse)
		http.Error(w, string(bytes), statusCode)
		return
	}

	jsonBytes, _ := json.Marshal(types.CourtCase{
		Cnj: courtCase.Cnj,
		Plaintiff: courtCase.Plaintiff,
		Defendant: courtCase.Defendant,
		CourtOfOrigin: courtCase.CourtOfOrigin,
		StartDate: courtCase.StartDate.Time,
	})
	w.Write(jsonBytes)
}

