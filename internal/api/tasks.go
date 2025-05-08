package api

import (
	"github.com/langowen/go_final_project/internal/db"
	"net/http"
)

const limit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(storage db.Storage, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		tasks, err = storage.SearchTasks(search, limit)
	} else {
		tasks, err = storage.Tasks(limit)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := TasksResp{
		Tasks: tasks,
	}
	if response.Tasks == nil {
		response.Tasks = []*db.Task{}
	}

	respondWithJSON(w, http.StatusOK, response)
}
