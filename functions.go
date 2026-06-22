package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getFileserverHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`<html>
  	<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  	</body>
	</html>`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetFileserverHits(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		log.Printf("Resetting fileserver hits is not allowed in %s environment\n", cfg.platform)
		w.WriteHeader(403)
		w.Write([]byte("Resetting fileserver hits is not allowed in this environment"))
		return
	}
	log.Printf("Resetting fileserver hits to 0\n")
	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.DeleteUser(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %s\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Failed to reset users"))
		return
	}
	w.WriteHeader(200)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnVals struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, returnVals{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON %s\n", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		respondWithError(w, 500, "Failed decoding")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	cleanBody := cleanWords(params.Body)
	respondWithJSON(w, 200, returnVals{Cleaned_body: cleanBody})

}

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
