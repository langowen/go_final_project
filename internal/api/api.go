package api

import (
	"encoding/json"
	"github.com/langowen/go_final_project/internal/db"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Init(todo *http.ServeMux, storage db.Storage) {
	todo.HandleFunc("/api/nextdate", nextDateHandler)
	todo.HandleFunc("/api/task", taskHandler(storage))
	todo.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasksHandler(storage, w, r)
	})
	todo.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		DoneTaskHandler(storage, w, r)
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
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
