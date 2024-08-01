package endpoints

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Healthcheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Healthy"))
}
