package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sondrefjellving/chirpy/internal/database"
)


type apiConfig struct {
	fileServerHits int
	db *database.DB
}

func main() {
	path := "database.json"
	db, err := database.NewDB(path)
	if err != nil {
		fmt.Println("Error creating db")
		return
	}
	mux := http.NewServeMux()
	cfg := apiConfig{
		fileServerHits: 0,
		db: db,
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", cfg.middlewareMetricsInc(handler))
	mux.Handle("/assets", http.FileServer(http.Dir(".logo.png")))
	mux.HandleFunc("GET /api/healthz", handlerReadinessGet)
	mux.HandleFunc("/api/reset", cfg.handleReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpPost)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpGet)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerServerHitsGet)
	
	server := &http.Server{
		Addr:			":8080",
		Handler:		mux,
	}

	fmt.Println("starting server...")
	log.Fatal(server.ListenAndServe())
}
