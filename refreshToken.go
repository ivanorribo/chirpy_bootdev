package main

import (
	"net/http"
	"time"

	"github.com/ivanorribo/chirpy_bootdev/internal/auth"
)

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}

	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid Authorization header")
		return
	}

	userID, err := cfg.dbQueries.CheckUserFromRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, 401, "Invalid refresh token")
		return
	}

	token, err := auth.MakeJWT(userID.ID, cfg.secretKey, time.Minute*60)
	if err != nil {
		respondWithError(w, 500, "Failed to create JWT")
		return
	}

	respondWithJSON(w, 200, returnVals{Token: token})
}
