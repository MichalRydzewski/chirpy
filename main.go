package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()
	cfg := apiConfig{}
	fileServer := http.FileServer(http.Dir(filepathRoot))

	srv := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(fileServer)))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("POST /admin/reset", cfg.handleReset)
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)

	fmt.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
