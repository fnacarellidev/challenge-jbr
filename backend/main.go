package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type CourtCase struct {
	Cnj           string    `json:"cnj"`
	Plaintiff     string    `json:"plaintiff"`
	Defendant     string    `json:"defendant"`
	CourtOfOrigin string    `json:"courtOfOrigin"`
	StartDate     time.Time `json:"startDate"`
}

// todo think about case where the court has already been registered
func RegisterCourtCase(w http.ResponseWriter, r *http.Request) {
	var courtCase CourtCase

	err := json.NewDecoder(r.Body).Decode(&courtCase)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Received info: %s", courtCase)
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to db")
	}
	defer conn.Close(context.Background())

	http.HandleFunc("/register_court_case", RegisterCourtCase)
	http.ListenAndServe(":8081", nil)
}
