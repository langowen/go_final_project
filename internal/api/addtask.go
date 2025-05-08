package api

import (
	"encoding/json"
	"github.com/langowen/go_final_project/internal/db"
	"net/http"
	"time"
)

func AddTaskHandler(storage db.Storage, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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

	//formatedDate := now.Format("20060102")
	//
	//nowFormated, err := time.Parse("20060102", formatedDate)
	//if err != nil {
	//	return err
	//}

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
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format("20060102")
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}

	return nil
}

// respondWithSuccess отправляет JSON с ID созданной задачи
func respondWithSuccess(w http.ResponseWriter, code int, id string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct {
		ID string `json:"id"`
	}{ID: id})
}
