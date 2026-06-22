package main

import "net/http"

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnVals struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, returnVals{Error: msg})
}
