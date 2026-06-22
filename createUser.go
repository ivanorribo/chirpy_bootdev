package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		respondWithError(w, 500, "Failed decoding")
		return
	}
	NewUser, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %s\n", err)
		respondWithError(w, 500, "Failed creating user")
		return
	}
	respondWithJSON(w, 201, User{ID: NewUser.ID, CreatedAt: NewUser.CreatedAt, UpdatedAt: NewUser.UpdatedAt, Email: NewUser.Email})
}
