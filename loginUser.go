package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ivanorribo/chirpy_bootdev/internal/auth"
	"github.com/ivanorribo/chirpy_bootdev/internal/database"
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

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Minute*60) // Token valid for 60 minutes
	if err != nil {
		log.Printf("Error creating JWT: %s\n", err)
		respondWithError(w, 500, "Failed to create JWT")
		return
	}

	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(((24 * time.Hour) * 60)), // Refresh token valid for 60 days
	})
	if err != nil {
		log.Printf("Error creating refresh token: %s\n", err)
		respondWithError(w, 500, "Failed to create refresh token")
		return
	}
	respondWithJSON(w, 200, User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, Token: token, RefreshToken: refreshToken})
}
