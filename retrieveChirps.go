package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) retrieveChirps(w http.ResponseWriter, r *http.Request) {

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	s := r.URL.Query().Get("author_id")
	if s != "" {
		authorID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, 400, "Invalid author ID")
			return
		}

		chirps, err := cfg.dbQueries.GetChirpsByAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, 500, "Failed retrieving chirps by author")
			return
		}
		sortDirection := r.URL.Query().Get("sort")
		sort.Slice(chirps, func(i, j int) bool {
			if sortDirection == "desc" {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			}
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})

		var returnChirps []returnVals
		for _, chirp := range chirps {
			returnChirps = append(returnChirps, returnVals{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}

		respondWithJSON(w, 200, returnChirps)
		return
	}

	chirps, err := cfg.dbQueries.RetrieveAll(r.Context())
	if err != nil {
		respondWithError(w, 500, "Failed retrieving chirps")
		return
	}

	sortDirection := r.URL.Query().Get("sort")
	sort.Slice(chirps, func(i, j int) bool {
		if sortDirection == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	var returnChirps []returnVals
	for _, chirp := range chirps {
		returnChirps = append(returnChirps, returnVals{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, 200, returnChirps)
}

func (cfg *apiConfig) retrieveChirpByID(w http.ResponseWriter, r *http.Request) {

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	parsedUUID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}

	chirp, error := cfg.dbQueries.RetrieveOneChirp(r.Context(), parsedUUID)
	if error != nil {
		respondWithError(w, 404, "Failed retrieving chirp")
		return
	}

	respondWithJSON(w, 200, returnVals{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
