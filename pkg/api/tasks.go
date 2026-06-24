package api

import (
	"net/http"
	"time"

	"go_final_pablo/pkg/db"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	var (
		tasks []*db.Task
		err   error
	)

	if search != "" {
		t, parseErr := time.Parse("02.01.2006", search)
		if parseErr == nil {
			tasks, err = db.TasksByDate(t.Format(dateFormat), 50)
		} else {
			tasks, err = db.TasksSearch(search, 50)
		}
	} else {
		tasks, err = db.Tasks(50)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tasksResp{Tasks: tasks})
}
