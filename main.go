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
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", cfg.middlewareMetricsInc(handler))
	mux.Handle("/assets", http.FileServer(http.Dir(".logo.png")))
	mux.HandleFunc("GET /api/healthz", handlerReadinessGet)
	mux.HandleFunc("GET /api/metrics", cfg.handlerServerHitsGet)
	mux.HandleFunc("/api/reset", cfg.handleReset)

	server := &http.Server{
		Addr:			":8080",
		Handler:		mux,
	}

	fmt.Println("starting server...")
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) handlerServerHitsGet(w http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf("Hits: %v", cfg.fileServerHits)
	w.Header().Add("Content-Type", "text-plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

