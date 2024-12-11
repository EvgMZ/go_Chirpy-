package handler

import (
	"chirpy/internal/database"
	"chirpy/reponse"
	"fmt"
	"net/http"
	"sync/atomic"
	"text/template"
)

type ApiConfig struct {
	Db             *database.Queries
	FileserverHits atomic.Int32
	Platform       string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем текущее значение счётчика
	hits := cfg.FileserverHits.Load()
	// Возвращаем количество запросов
	fmt.Fprintf(w, "Hits: %d\n", hits)
}
func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	// // Сбрасываем атомарный счётчик запросов
	// cfg.fileserverHits.Store(0)
	// // Возвращаем подтверждение сброса
	// fmt.Fprintln(w, "Hits counter reset to 0")
	type response struct {
		Body string `json: "body"`
	}
	if cfg.Platform != "dev" {
		reponse.RespondWithError(w, 403, "Forbidden")
		return
	}
	_, err := cfg.Db.DeleteAllUser(r.Context())
	if err != nil {
		reponse.RespondWithError(w, 500, "user don't delete")
	}
	reponse.RespondWithJSON(w, 200, response{Body: "ok"})

}

func (cfg *ApiConfig) AdminMetricsHanlder(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index_admin.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	data := struct {
		Hits int32
	}{
		Hits: cfg.FileserverHits.Load(),
	}

	// Рендерим шаблон с данными
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	// http.ServeFile(w, r, "index_admin.html")
}
