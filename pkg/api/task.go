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
		respondWithError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	task, err := storage.Task(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			respondWithError(w, http.StatusNotFound, "задача не найдена")
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
		respondWithError(w, http.StatusBadRequest, "неверный JSON формат")
		return
	}

	if task.Title == "" {
		respondWithError(w, http.StatusBadRequest, "отсутствует заголовок задачи")
		return
	}

	err := checkDate(&task)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = storage.UpdateTask(&task)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "не удалось обновить задачу: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, TaskResp{})
}

func DoneTaskHandler(storage db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			respondWithError(w, http.StatusBadRequest, "не указан идентификатор")
			return
		}

		task, err := storage.Task(id)
		if err != nil {
			if errors.Is(err, db.ErrTaskNotFound) {
				respondWithError(w, http.StatusNotFound, "задача не найдена")
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
			respondWithError(w, http.StatusInternalServerError, "не удалось обновить задачу: "+err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, struct{}{})
	}
}

func DelTaskHandler(storage db.Storage, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	err := storage.Delete(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
