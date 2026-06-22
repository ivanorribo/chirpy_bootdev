package main

import (
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
