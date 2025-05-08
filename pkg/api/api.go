package api

import (
	"encoding/json"
	"github.com/langowen/go_final_project/pkg/auth"
	"github.com/langowen/go_final_project/pkg/config"
	"github.com/langowen/go_final_project/pkg/db"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Init(mux *http.ServeMux, cfg *config.Config, storage db.Storage) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/task", auth.Middleware(cfg, taskHandler(storage)))
	mux.HandleFunc("/api/tasks", auth.Middleware(cfg, tasksHandler(storage)))
	mux.HandleFunc("/api/task/done", auth.Middleware(cfg, DoneTaskHandler(storage)))
	mux.HandleFunc("/api/signin", func(w http.ResponseWriter, r *http.Request) {
		SignInHandler(cfg, w, r)
	})
}

func taskHandler(storage db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			AddTaskHandler(storage, w, r)
		case http.MethodGet:
			GetTaskHandler(storage, w, r)
		case http.MethodPut:
			PutTaskHandler(storage, w, r)
		case http.MethodDelete:
			DelTaskHandler(storage, w, r)
		default:
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(ErrorResponse{Error: message})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
