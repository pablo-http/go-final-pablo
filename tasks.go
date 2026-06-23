package main

import (
	"fmt"
	"net/http"
	"time"
)

// Tasks возвращает список ближайших задач из БД.
func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения задач: %w", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка сканирования задачи: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// TasksSearch возвращает задачи по строке поиска или дате.
func TasksSearch(search string, limit int) ([]*Task, error) {
	// Проверяем, является ли search датой в формате 02.01.2006
	t, err := time.Parse("02.01.2006", search)
	if err == nil {
		date := t.Format(dateFormat)
		rows, err := db.Query(
			`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`,
			date, limit,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка поиска по дате: %w", err)
		}
		defer rows.Close()

		tasks := make([]*Task, 0)
		for rows.Next() {
			task := &Task{}
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}
		return tasks, rows.Err()
	}

	// Поиск по подстроке в title или comment
	like := "%" + search + "%"
	rows, err := db.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
		like, like, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска задач: %w", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// tasksHandler обрабатывает GET /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	type TasksResp struct {
		Tasks []*Task `json:"tasks"`
	}

	search := r.FormValue("search")

	var (
		tasks []*Task
		err   error
	)

	if search != "" {
		tasks, err = TasksSearch(search, 50)
	} else {
		tasks, err = Tasks(50)
	}

	if err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
