package main

import (
	handler "chirpy/hanlder"
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func healthzFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Errorf("error connect to db %v", err)
	}
	dbQueries := database.New(db)
	apiCfg := handler.ApiConfig{
		Db:       dbQueries,
		Platform: platform,
	}
	// apiCfg := &hanlder.apiConfig{
	// 	Db:       dbQueries,
	// 	platform: platform,
	// }
	servMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: servMux,
	}
	// servMux.Handle("/app/", apiCfg.middlewareMetricsInc(server.Handler))
	servMux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	servMux.HandleFunc("/admin/metrics", apiCfg.AdminMetricsHanlder)
	servMux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	servMux.HandleFunc("POST /admin/reset", apiCfg.ResetHandler)
	servMux.HandleFunc("GET /api/healthz", healthzFunc)
	servMux.HandleFunc("POST /api/users", apiCfg.CreateUser)
	servMux.HandleFunc("POST /api/validate_chirp", handler.NewValidateChripHandler)
	// servMux.HandleFunc("POST /api/validate_chirp", validateChripHandler)

	server.ListenAndServe()
}
