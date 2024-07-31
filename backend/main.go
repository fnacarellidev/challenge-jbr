package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fnacarellidev/challenge-jbr/backend/.sqlcbuild/pgquery"
	"github.com/fnacarellidev/challenge-jbr/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
)

type Api struct {
	Db *pgx.Conn
}

type CaseUpdate struct {
	UpdateDate    time.Time `json:"update_date"`
	UpdateDetails string    `json:"update_details"`
}

// todo think about case where the court has already been registered
func (api *Api) RegisterCourtCase(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var courtCase types.CourtCase

	err := json.NewDecoder(r.Body).Decode(&courtCase)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	sqlc := pgquery.New(api.Db)
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

func (api *Api) FetchCourtCase(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cnjLookup := ps.ByName("cnj")
	sqlc := pgquery.New(api.Db)

	courtCase, error := sqlc.GetCourtCase(context.Background(), cnjLookup)
	if error != nil {
		http.Error(w, "There's no such court case", http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(types.CourtCase{
		Cnj: courtCase.Cnj,
		Plaintiff: courtCase.Plaintiff,
		Defendant: courtCase.Defendant,
		CourtOfOrigin: courtCase.CourtOfOrigin,
		StartDate: courtCase.StartDate.Time,
	})
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonBytes)
}

func (api *Api) FetchUpdatesFromCase(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cnjLookup := ps.ByName("cnj")
	sqlc := pgquery.New(api.Db)

	rows, err := sqlc.GetCaseUpdates(context.Background(), cnjLookup)
	if err != nil {
		log.Println("err:", err)
		return
	}

	var CaseUpdates []CaseUpdate
	for _, row := range(rows) {
		CaseUpdate := CaseUpdate{
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

func (api *Api) Healthcheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Healthy"))
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to db")
	}
	api := &Api{Db: conn}
	defer conn.Close(context.Background())

	router := httprouter.New()
	router.POST("/register_court_case", api.RegisterCourtCase)
	router.GET("/fetch_court_case/:cnj", api.FetchCourtCase)
	router.GET("/fetch_updates_from_case/:cnj", api.FetchUpdatesFromCase)
	router.GET("/healthcheck", api.Healthcheck)
	http.ListenAndServe(":8081", router)
}
