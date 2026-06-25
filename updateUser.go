package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ivanorribo/chirpy_bootdev/internal/auth"
	"github.com/ivanorribo/chirpy_bootdev/internal/database"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type returnVals struct {
		UserID    uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		respondWithError(w, 400, "Invalid request body")
		return
	}

	userID, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error retrieving user ID from token: %s\n", err)
		respondWithError(w, 401, "Unauthorized")
		return
	}

	validUserID, err := auth.ValidateJWT(userID, cfg.secretKey)
	if err != nil {
		log.Printf("Error validating JWT: %s\n", err)
		respondWithError(w, 401, "Unauthorized")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s\n", err)
		respondWithError(w, 500, "Failed to hash password")
		return
	}

	updatedUser, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             validUserID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error updating user: %s\n", err)
		respondWithError(w, 500, "Failed to update user")
		return
	}

	respondWithJSON(w, 200, returnVals{
		UserID:    updatedUser.ID,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	})
}
