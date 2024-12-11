package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func healthzFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// cfg.fileserverHits.Add(1)
	// return next
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем текущее значение счётчика
	hits := cfg.fileserverHits.Load()
	// Возвращаем количество запросов
	fmt.Fprintf(w, "Hits: %d\n", hits)
}
func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	// Сбрасываем атомарный счётчик запросов
	cfg.fileserverHits.Store(0)
	// Возвращаем подтверждение сброса
	fmt.Fprintln(w, "Hits counter reset to 0")
}
func main() {
	apiCfg := &apiConfig{}
	servMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: servMux,
	}
	// servMux.Handle("/app/", apiCfg.middlewareMetricsInc(server.Handler))
	servMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	servMux.Handle("GET /admin/metrics", apiCfg.middlewareMetricsInc(http.FileServer("./index_admin.html")))
	servMux.HandleFunc("GET /api/healthz", healthzFunc)
	servMux.HandleFunc("GET /api/metrics", apiCfg.metricsHandler)
	servMux.HandleFunc("POST /api/reset", apiCfg.resetHandler)

	server.ListenAndServe()
}
