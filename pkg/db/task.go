package db

import (
	"fmt"
	"log"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("AddTask error: %v", err)
		return 0, fmt.Errorf("внутренняя ошибка")
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	t := &Task{}
	err := DB.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id,
	).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}
	return t, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		log.Printf("UpdateTask error: %v", err)
		return fmt.Errorf("внутренняя ошибка")
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("внутренняя ошибка")
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		log.Printf("DeleteTask error: %v", err)
		return fmt.Errorf("внутренняя ошибка")
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("внутренняя ошибка")
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateDate(id string, next string) error {
	_, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		log.Printf("UpdateDate error: %v", err)
		return fmt.Errorf("внутренняя ошибка")
	}
	return nil
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`,
		limit,
	)
	if err != nil {
		log.Printf("Tasks error: %v", err)
		return nil, fmt.Errorf("внутренняя ошибка")
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("внутренняя ошибка")
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func TasksSearch(search string, limit int) ([]*Task, error) {
	like := "%" + search + "%"
	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
		like, like, limit,
	)
	if err != nil {
		log.Printf("TasksSearch error: %v", err)
		return nil, fmt.Errorf("внутренняя ошибка")
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("внутренняя ошибка")
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func TasksByDate(date string, limit int) ([]*Task, error) {
	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`,
		date, limit,
	)
	if err != nil {
		log.Printf("TasksByDate error: %v", err)
		return nil, fmt.Errorf("внутренняя ошибка")
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("внутренняя ошибка")
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}
