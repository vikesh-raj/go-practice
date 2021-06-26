package server

import (
	"encoding/json"
	"net/http"
)

func writeResponse(w http.ResponseWriter, r interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(r)
}

func writeErrorResponse(w http.ResponseWriter, message string, code int) {
	res := BasicResponse{
		Status:       "FAIL",
		ErrorMessage: message,
	}
	writeResponse(w, res, code)
}
