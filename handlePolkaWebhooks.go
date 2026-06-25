package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeChirpyToRed(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	param := &paramaters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(param)
	if err != nil {
		respondWithError(w, 400, "Invalid request body")
		return
	}

	if param.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent) // 204 No Content
		return
	}

	_, err = cfg.dbQueries.UpgradeToChirpyRed(r.Context(), param.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "Failed to upgrade user")
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
