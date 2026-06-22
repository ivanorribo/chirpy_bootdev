package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ivanorribo/chirpy_bootdev/internal/database"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		respondWithError(w, 500, "Failed decoding")
		return
	}
	// first validate the chirp and clean it
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	params.Body = cleanWords(params.Body)
	Chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s\n", err)
		respondWithError(w, 500, "Failed creating chirp")
		return
	}
	respondWithJSON(w, 201, returnVals{
		ID:        Chirp.ID,
		CreatedAt: Chirp.CreatedAt,
		UpdatedAt: Chirp.UpdatedAt,
		Body:      Chirp.Body,
		UserID:    Chirp.UserID,
	})
}
