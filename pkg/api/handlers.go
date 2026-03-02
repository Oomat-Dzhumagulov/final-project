package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Oomat-Dzhumagulov/final-project/pkg/scheduler"
)

const DateFormat = "20060102"

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "метод не подерживается", http.StatusMethodNotAllowed)
	}
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var (
		now time.Time
		err error
	)

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка формата now: %v", err), http.StatusBadRequest)
			return
		}

		nextDateStr, err := scheduler.NextDate(now, dateStr, repeatStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка вычисление следущей даты: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(nextDateStr))
	}
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"error": msg,
	})
}
