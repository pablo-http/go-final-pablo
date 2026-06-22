package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// checkDate проверяет и корректирует дату задачи.
func checkDate(task *Task) error {
	now := time.Now()

	// Если дата не указана — подставляем сегодня
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
		return nil
	}

	// Проверяем корректность формата даты
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	// Если правило повторения указано — проверяем его и получаем следующую дату
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// Если дата задачи меньше сегодняшней
	if !t.After(truncateToDay(now)) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

// truncateToDay обнуляет время до начала суток.
func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// addTaskHandler обрабатывает POST /api/task.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "ошибка десериализации JSON: "+err.Error())
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

	id, err := AddTask(&task)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{"id": id})
}
