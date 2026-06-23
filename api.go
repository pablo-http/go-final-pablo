package main

import (
	"encoding/json"
	"net/http"
)

// writeJSON сериализует data в JSON и отправляет ответ
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// writeError отправляет JSON с полем error
func writeError(w http.ResponseWriter, msg string) {
	writeJSON(w, map[string]string{"error": msg})
}

// taskHandler маршрутизирует запросы /api/task по HTTP-методу
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeError(w, "метод не поддерживается")
	}
}

// getTaskHandler обрабатывает GET /api/task?id=...
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, task)
}

// updateTaskHandler обрабатывает PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "ошибка десериализации JSON: "+err.Error())
		return
	}

	if task.ID == "" {
		writeError(w, "не указан идентификатор")
		return
	}

	if task.Title == "" {
		writeError(w, "не указан заголовок задачи")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	if err := UpdateTask(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{})
}

// deleteTaskHandler обрабатывает DELETE /api/task?id=...
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор")
		return
	}

	if err := DeleteTask(id); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{})
}
