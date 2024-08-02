package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/fnacarellidev/challenge-jbr/backend/endpoints"
)

func main() {
	router := httprouter.New()
	router.POST("/register_court_case", endpoints.RegisterCourtCase)
	router.GET("/fetch_court_case/:cnj", endpoints.FetchCourtCase)
	router.GET("/healthcheck", endpoints.Healthcheck)
	http.ListenAndServe(":8081", router)
}
