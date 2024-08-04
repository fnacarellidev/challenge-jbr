package endpoints

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/fnacarellidev/challenge-jbr/backend/.sqlcbuild/pgquery"
	"github.com/fnacarellidev/challenge-jbr/backend/endpoints/utils"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
)

func RegisterCourtCase(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	defer conn.Close(context.Background())
	if err != nil {
		log.Fatal("Failed to connect to db")
	}

	var courtCase types.CourtCase
	err = json.NewDecoder(r.Body).Decode(&courtCase)
	if err != nil {
		utils.SendError(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := courtCase.Validate(); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlc := pgquery.New(conn)
	_, err = sqlc.InsertCourtCase(context.Background(), pgquery.InsertCourtCaseParams{
		Cnj: courtCase.Cnj,
		Plaintiff: courtCase.Plaintiff,
		Defendant: courtCase.Defendant,
		CourtOfOrigin: courtCase.CourtOfOrigin,
		StartDate: pgtype.Date{
			Time: courtCase.StartDate,
			Valid: true,
		},
	})
	if err != nil {
		utils.SendError(w, "case already exists", http.StatusConflict)
	}

	for _, update := range courtCase.Updates {
		sqlc.InsertCaseUpdate(context.Background(), pgquery.InsertCaseUpdateParams{
			Cnj: courtCase.Cnj,
			UpdateDate: pgtype.Timestamptz{
				Time: update.UpdateDate,
				Valid: true,
			},
			UpdateDetails: update.UpdateDetails,
		})
	}
}
