package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/ivanorribo/chirpy_bootdev/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	handlerCfg     http.Handler
	dbQueries      *database.Queries
	platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM environment variable is not set")
	}
	apiCfg := &apiConfig{platform: platform}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	apiCfg.dbQueries = dbQueries
	const port = ":8080"
	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))

	mux.HandleFunc("GET /api/healthz", handlerReadiness) // add readiness handler using function from functions.go
	mux.HandleFunc("GET /admin/metrics", apiCfg.getFileserverHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetFileserverHits)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	stripHandler := http.StripPrefix("/app/", handler) // strip the /app prefix so we can differentiate between the /app/ path and the root path
	mux.Handle("/app", http.RedirectHandler("/app/", http.StatusPermanentRedirect))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(stripHandler))

	Server := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(Server.ListenAndServe())
}
