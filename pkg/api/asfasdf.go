package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func beforeOrEqualDate(date, now time.Time) bool {
	y1, m1, d1 := now.Date()
	y2, m2, d2 := date.Date()
	return y1 < y2 || (y1 == y2 && m1 < m2) || (y1 == y2 && m1 == m2 && d1 < d2)
}

func checkRules(dstart string, repeat string) (time.Time, string, int, error) {
	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		return time.Time{}, "", -1, errors.New("некорректный формат dstart")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return time.Time{}, "", -1, errors.New("пустое правило повторения")
	}

	rule := parts[0]
	var interval int

	switch rule {
	case "d":
		if len(parts) != 2 {
			return time.Time{}, "", -1, errors.New("неверный формат правила 'd'")
		}
		interval, err = strconv.Atoi(parts[1])
		if err != nil || interval <= 0 || interval > 400 {
			return time.Time{}, "", -1, errors.New("недопустимый интервал для правила 'd'")
		}

	case "y":
		if len(parts) != 1 {
			return time.Time{}, "", -1, errors.New("неверный формат правила 'y'")
		}

	default:
		return time.Time{}, "", -1, fmt.Errorf("неподдерживаемый формат правила '%s'", rule)
	}

	return startDate, rule, interval, nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	nowWithoutTime := now
	startDate, rule, interval, err := checkRules(dstart, repeat)
	if err != nil {
		return "", err
	}
	currentDate := startDate
	for {
		switch rule {
		case "d":
			currentDate = currentDate.AddDate(0, 0, interval)
			if beforeOrEqualDate(currentDate, nowWithoutTime) {
				return currentDate.Format("20060102"), nil
			}
		case "y":
			currentDate = currentDate.AddDate(1, 0, 0)
			if beforeOrEqualDate(currentDate, nowWithoutTime) {
				return currentDate.Format("20060102"), nil
			}
		}
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	if dateParam == "" || repeatParam == "" {
		http.Error(w, "Отсутствуют обязательные параметры 'date' или 'repeat'", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowParam == "" {
		now = time.Now().UTC()
	} else {
		parsedNow, err := time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, "Некорректный формат параметра 'now'", http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка: %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, nextDate)

}
