package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ivanorribo/chirpy_bootdev/internal/auth"
	"github.com/ivanorribo/chirpy_bootdev/internal/database"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		respondWithError(w, 500, "Failed decoding")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s\n", err)
		respondWithError(w, 500, "Failed hashing password")
		return
	}
	NewUser, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %s\n", err)
		respondWithError(w, 500, "Failed creating user")
		return
	}
	respondWithJSON(w, 201, User{ID: NewUser.ID, CreatedAt: NewUser.CreatedAt, UpdatedAt: NewUser.UpdatedAt, Email: NewUser.Email})
}
