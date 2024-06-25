package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sondrefjellving/chirpy/internal/database"
)


type apiConfig struct {
	fileServerHits int
	db *database.DB
	jwtSecret string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	path := "database.json"
	db, err := database.NewDB(path, dbg)
	if err != nil {
		log.Fatal("Error creating db")
		return
	}

	mux := http.NewServeMux()
	cfg := apiConfig{
		fileServerHits: 0,
		db: db,
		jwtSecret: jwtSecret,
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", cfg.middlewareMetricsInc(handler))
	mux.Handle("/assets", http.FileServer(http.Dir(".logo.png")))
	mux.HandleFunc("GET /api/healthz", handlerReadinessGet)
	mux.HandleFunc("/api/reset", cfg.handleReset)

	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpPost)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpGet)
	mux.HandleFunc("GET /api/chirps/{chirpId}", cfg.handlerChirpGetById)

	mux.HandleFunc("POST /api/users", cfg.handlerUserPost)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerServerHitsGet)

	mux.HandleFunc("POST /api/login", cfg.handlerLoginPost)

	mux.HandleFunc("PUT /api/users", cfg.handlerUserPut)
	
	server := &http.Server{
		Addr:			":8080",
		Handler:		mux,
	}

	fmt.Println("starting server...")
	log.Fatal(server.ListenAndServe())
}
