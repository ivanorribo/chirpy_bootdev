package main

import (
	"net/http"

	"github.com/ivanorribo/chirpy_bootdev/internal/auth"
)

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid Authorization header")
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, 500, "Failed to revoke refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
