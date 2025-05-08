package api

import (
	"encoding/json"
	"errors"
	"github.com/langowen/go_final_project/pkg/db"
	"net/http"
	"time"
)

type TaskResp struct {
	Task *db.Task `json:"task"`
}

func GetTaskHandler(storage db.Storage, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, http.StatusMethodNotAllowed, "Не указан идентификатор")
		return
	}

	task, err := storage.Task(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			respondWithError(w, http.StatusMethodNotAllowed, "Задача не найдена")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func PutTaskHandler(storage db.Storage, w http.ResponseWriter, r *http.Request) {
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

	err = storage.UpdateTask(&task)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed update task: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, TaskResp{})
}

func DoneTaskHandler(storage db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			respondWithError(w, http.StatusMethodNotAllowed, "Не указан идентификатор")
			return
		}

		task, err := storage.Task(id)
		if err != nil {
			if errors.Is(err, db.ErrTaskNotFound) {
				respondWithError(w, http.StatusMethodNotAllowed, "Задача не найдена")
				return
			}
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if task.Repeat == "" {
			err = storage.Delete(id)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondWithJSON(w, http.StatusOK, struct{}{})
			return
		}

		now := time.Now()

		task.Date, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = storage.UpdateTask(task)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed update task: "+err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, struct{}{})
	}
}

func DelTaskHandler(storage db.Storage, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, http.StatusMethodNotAllowed, "Не указан идентификатор")
		return
	}

	err := storage.Delete(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
