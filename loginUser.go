package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ivanorribo/chirpy_bootdev/internal/auth"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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
	if params.ExpiresInSeconds <= 0 || params.ExpiresInSeconds > 3600 {
		params.ExpiresInSeconds = 3600 // default to 1 hour if not provided or invalid
	}
	token, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		log.Printf("Error creating JWT: %s\n", err)
		respondWithError(w, 500, "Failed to create JWT")
		return
	}

	respondWithJSON(w, 200, User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, Token: token})
}
