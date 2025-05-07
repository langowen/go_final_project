package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Init(todo *http.ServeMux) {
	todo.HandleFunc("/api/nextdate", nextDateHandler)
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if nowStr == "" || date == "" || repeat == "" {
		respondWithError(w, http.StatusBadRequest, "missing required parameters: now, date or repeat")
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid now parameter: "+err.Error())
		return
	}

	result, err := NextDate(now, date, repeat)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
