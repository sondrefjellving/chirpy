package main

import (
	"fmt"
	"log"
	"net/http"
)


type apiConfig struct {
	fileServerHits int
}

func main() {
	mux := http.NewServeMux()
	cfg := apiConfig{
		fileServerHits: 0,
		chirps:	[]Chirp{},
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", cfg.middlewareMetricsInc(handler))
	mux.Handle("/assets", http.FileServer(http.Dir(".logo.png")))
	mux.HandleFunc("GET /api/healthz", handlerReadinessGet)
	mux.HandleFunc("/api/reset", cfg.handleReset)
	mux.HandleFunc("POST /api/chirps", handlerValidateChirpPost)
	mux.HandleFunc("GET /api/chirps", handlerGetChirps)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerServerHitsGet)
	
	server := &http.Server{
		Addr:			":8080",
		Handler:		mux,
	}

	fmt.Println("starting server...")
	log.Fatal(server.ListenAndServe())
}
