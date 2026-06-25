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

	type returnVals struct {
		UserID       string    `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		respondWithError(w, 400, "Invalid request body")
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

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Hour) // Token valid for 60 minutes
	if err != nil {
		log.Printf("Error creating JWT: %s\n", err)
		respondWithError(w, 500, "Failed to create JWT")
		return
	}

	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour), // Refresh token valid for 60 days
	})
	if err != nil {
		log.Printf("Error creating refresh token: %s\n", err)
		respondWithError(w, 500, "Failed to create refresh token")
		return
	}
	respondWithJSON(w, 200, returnVals{
		UserID:       user.ID.String(),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
