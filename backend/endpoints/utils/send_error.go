package utils

import (
	"encoding/json"
	"net/http"

	"github.com/fnacarellidev/challenge-jbr/types"
)

func SendError(w http.ResponseWriter, msg string, code int) {
	e := types.ErrResponse{
		Error: msg,
	}
	bytes, _ := json.Marshal(e)
	http.Error(w, string(bytes), code)
}
