package main

import (
	"fmt"
	"io"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir("./app")))
	apiCfg := &apiConfig{0}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fsHandler))
	mux.HandleFunc("/healthz", ready)
	mux.HandleFunc("/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("/reset", apiCfg.resetHits)
	corsMux := middlewareCors(mux)
	http.ListenAndServe(":8080", corsMux)
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
		fmt.Printf("Hits: %d\n", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Printf("metrics hit, hits = %d\n", cfg.fileserverHits)
	io.WriteString(w, fmt.Sprintf("Hits: %d", cfg.fileserverHits))
	return
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = 0
	fmt.Printf("fileserverHits reset to 0\n")
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
