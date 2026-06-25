package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ivanorribo/chirpy_bootdev/internal/auth"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {

	parsedUUID, err := uuid.Parse(r.PathValue("chirpID")) // Parse the chirpID from the URL path
	if err != nil {
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}

	token, err := auth.GetBearerToken(r.Header) // Get the JWT token from the Authorization header
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	validUserID, err := auth.ValidateJWT(token, cfg.secretKey) // Validate the JWT token
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	chirps, err := cfg.dbQueries.RetrieveOneChirp(r.Context(), parsedUUID) // Retrieve the chirp from the database
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	if chirps.UserID != validUserID {
		respondWithError(w, 403, "Forbidden: You can only delete your own chirps")
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), parsedUUID) // Delete the chirp from the database
	if err != nil {
		respondWithError(w, 500, "Failed deleting chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
