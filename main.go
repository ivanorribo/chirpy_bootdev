package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	handlerCfg     http.Handler
}

func main() {
	const port = ":8080"
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}
	handler := http.FileServer(http.Dir("."))

	mux.HandleFunc("GET /healthz", handlerReadiness) // add readiness handler using function from functions.go
	mux.HandleFunc("GET /metrics", apiCfg.getFileserverHits)
	mux.HandleFunc("POST /reset", apiCfg.resetFileserverHits)

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
