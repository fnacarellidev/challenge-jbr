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

func FetchUpdatesFromCase(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to db")
	}
	defer conn.Close(context.Background())

	cnjLookup := ps.ByName("cnj")
	sqlc := pgquery.New(conn)
	rows, err := sqlc.GetCaseUpdates(context.Background(), cnjLookup)
	if err != nil {
		log.Println("err:", err)
		return
	}

	var CaseUpdates []types.CaseUpdate
	for _, row := range(rows) {
		CaseUpdate := types.CaseUpdate{
			UpdateDate: row.UpdateDate.Time,
			UpdateDetails: row.UpdateDetails,
		}
		CaseUpdates = append(CaseUpdates, CaseUpdate)
	}

	jsonBytes, err := json.Marshal(CaseUpdates)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonBytes)
}
