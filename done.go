package main

import (
	"net/http"
	"time"
)

// doneTaskHandler обрабатывает POST /api/task/done?id=...
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор")
		return
	}

	task, err := GetTask(id)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	// Если правило повторения не указано — удаляем задачу
	if task.Repeat == "" {
		if err := DeleteTask(id); err != nil {
			writeError(w, err.Error())
			return
		}
		writeJSON(w, map[string]any{})
		return
	}

	// Для периодической задачи вычисляем следующую дату
	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	if err := UpdateDate(id, next); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{})
}
