package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ivanorribo/chirpy_bootdev/internal/auth"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		respondWithError(w, 500, "Failed decoding")
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error retrieving user: %s\n", err)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		log.Printf("Error checking password hash: %s\n", err)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	respondWithJSON(w, 200, User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email})
}
