package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Oomat-Dzhumagulov/final-project/pkg/db"
	"github.com/Oomat-Dzhumagulov/final-project/pkg/scheduler"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPut:
		putTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeError(w, "Метод не подерживается", http.StatusMethodNotAllowed)
	}
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	var (
		tasks []*db.Task
		err   error
	)

	search := r.FormValue("search")
	if search != "" {
		tasks, err = db.TasksBySerch(search, 50)
	} else {
		tasks, err = db.Tasks(50)
	}
	if err != nil {
		writeError(w, "ошибка получания задач", http.StatusInternalServerError)
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "ошибка дисериализации", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "ошибка заголовка", http.StatusBadRequest)
		return
	}

	if err := correctTaskDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "ошибка добавления задачи", http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]any{
		"id": strconv.Itoa(int(id)),
	})
}

func correctTaskDate(task *db.Task) error {
	now := time.Now()
	nowStr := now.Format("20060102")

	if task.Date == "" {
		task.Date = nowStr
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("ошибка парсинга: %w", err)
	}

	var next string

	if task.Repeat != "" {
		next, err = scheduler.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("ошибка repeat: %w", err)
		}
	}

	if scheduler.IsAfter(now, t) {
		if task.Repeat == "" {
			task.Date = nowStr
		} else {
			task.Date = next
		}
	}

	return nil
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "задача не найдена", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJson(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "шибка десериализации", http.StatusBadRequest)
		return
	}
	if task.ID == "" {
		writeError(w, "задача не найдена", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeError(w, "заголовок не найден", http.StatusBadRequest)
		return
	}

	if err := correctTaskDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJson(w, map[string]any{})
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "мтод не подержива6тся", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeError(w, "ID не указан", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if err = db.DeleteTask(id); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		next, err := scheduler.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = db.UpdateDate(id, next); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	writeJson(w, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "ID не указан", http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJson(w, map[string]any{})
}
