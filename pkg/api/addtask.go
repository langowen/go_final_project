package api

import (
	"encoding/json"
	"github.com/langowen/go_final_project/pkg/db"
	"net/http"
	"time"
)

type SuccessResponse struct {
	ID string `json:"id"`
}

func AddTaskHandler(storage db.Storage, w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if task.Title == "" {
		respondWithError(w, http.StatusBadRequest, "Title is required")
		return
	}

	err := checkDate(&task)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := storage.AddTask(&task)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add task: "+err.Error())
		return
	}

	respondWithSuccess(w, http.StatusCreated, id)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}

	return nil
}

func respondWithSuccess(w http.ResponseWriter, code int, id string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(SuccessResponse{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
