package main

import (
	"encoding/json"
	"net/http"
)

// writeJSON сериализует data в JSON и отправляет ответ.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// writeError отправляет JSON с полем error.
func writeError(w http.ResponseWriter, msg string) {
	writeJSON(w, map[string]string{"error": msg})
}

// taskHandler маршрутизирует запросы /api/task по HTTP-методу.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		writeError(w, "метод не поддерживается")
	}
}
