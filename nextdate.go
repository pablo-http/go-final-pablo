package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

// NextDate вычисляет следующую дату выполнения задачи
// now — дата, относительно которой ищем следующую
// dstart — исходная дата задачи в формате 20060102
// repeat — правило повторения
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата: %w", err)
	}

	parts := strings.SplitN(repeat, " ", 2)
	switch parts[0] {

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}

	case "d":
		if len(parts) < 2 {
			return "", errors.New("правило d: не указан интервал")
		}
		interval, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || interval < 1 || interval > 400 {
			return "", fmt.Errorf("правило d: недопустимый интервал %q", parts[1])
		}
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				break
			}
		}

	case "w":
		if len(parts) < 2 {
			return "", errors.New("правило w: не указаны дни недели")
		}
		var weekdays [8]bool // индексы 1–7: пн–вс
		for _, s := range strings.Split(parts[1], ",") {
			d, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil || d < 1 || d > 7 {
				return "", fmt.Errorf("правило w: недопустимый день недели %q", s)
			}
			weekdays[d] = true
		}
		for {
			date = date.AddDate(0, 0, 1)
			wd := int(date.Weekday()) // 0=вс, 1=пн ... 6=сб
			if wd == 0 {
				wd = 7
			}
			if weekdays[wd] && date.After(now) {
				break
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", errors.New("правило m: не указаны дни")
		}
		// разбивка на дни и  месяцы
		mparts := strings.SplitN(parts[1], " ", 2)
		var days [32]bool // индексы 1–31, -1 и -2  отдельно
		var lastDay bool  // -1: последний день месяца
		var preLast bool  // -2: предпоследний день месяца
		var months [13]bool
		hasMonths := false

		for _, s := range strings.Split(mparts[0], ",") {
			d, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil || d < -2 || d == 0 || d > 31 {
				return "", fmt.Errorf("правило m: недопустимый день %q", s)
			}
			if d == -1 {
				lastDay = true
			} else if d == -2 {
				preLast = true
			} else {
				days[d] = true
			}
		}

		if len(mparts) == 2 {
			hasMonths = true
			for _, s := range strings.Split(mparts[1], ",") {
				mo, err := strconv.Atoi(strings.TrimSpace(s))
				if err != nil || mo < 1 || mo > 12 {
					return "", fmt.Errorf("правило m: недопустимый месяц %q", s)
				}
				months[mo] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			if !date.After(now) {
				continue
			}
			if hasMonths && !months[int(date.Month())] {
				continue
			}
			// последний и предпоследний день месяца
			last := lastDayOfMonth(date)
			if lastDay && date.Day() == last {
				break
			}
			if preLast && date.Day() == last-1 {
				break
			}
			if date.Day() <= 31 && days[date.Day()] {
				break
			}
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %q", repeat)
	}

	return date.Format(dateFormat), nil
}

// lastDayOfMonth возвращает номер последнего дня месяца для заданной даты
func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

// nextDateHandler обрабатывает GET /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "некорректный параметр now", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, next)
}
