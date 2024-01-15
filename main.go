package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	r := chi.NewRouter()
	r.Use(middlewareCors)
	apiCfg := &apiConfig{0}
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app"))))

	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)
	r.Get("/healthz", ready)
	r.Get("/metrics", apiCfg.hitsHandler)
	r.HandleFunc("/reset", apiCfg.resetHits)
	http.ListenAndServe(":8080", r)
}

func ready(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
	return
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, fmt.Sprintf("Hits: %d", cfg.fileserverHits))
	return
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	return
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
