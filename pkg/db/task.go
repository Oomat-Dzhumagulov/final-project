package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	tasks := make([]*Task, 0)

	query := `SELECT id, date, title, comment, repeat
			  FROM scheduler ORDER BY date LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return tasks, fmt.Errorf("ошибка запросак к БД: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var task Task

		var id int

		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, fmt.Errorf("ошибка чтения БД: %w", err)
		}

		task.ID = strconv.Itoa(int(id))

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return tasks, fmt.Errorf("ошибка перебора БД: %w", err)
	}
	return tasks, nil
}

func TasksBySerch(search string, limit int) ([]*Task, error) {
	tasks := make([]*Task, 0)

	var rows *sql.Rows

	t, err := time.Parse("02.01.2006", search)
	if err == nil {
		date := t.Format("20060102")
		rows, err = DB.Query(
			`SELECT id, date, title, comment, repeat 
			 FROM scheduler 
			 WHERE date = ? 
			 LIMIT ?`,
			date, limit,
		)
		if err != nil {
			return tasks, fmt.Errorf("ошибка запроса к БД: %w", err)
		}
	} else {
		like := "%" + search + "%"

		rows, err = DB.Query(
			`SELECT id, date, title, comment, repeat 
			 FROM scheduler 
			 WHERE title LIKE ? OR comment LIKE ? 
			 ORDER BY date 
			 LIMIT ?`,
			like, like, limit,
		)
		if err != nil {
			return tasks, fmt.Errorf("ошибка запроса к БД: %w", err)
		}
	}

	defer rows.Close()
	for rows.Next() {
		var (
			task Task
			id   int64
		)
		err = rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, fmt.Errorf("ошибка чтения БД: %w", err)
		}

		task.ID = strconv.Itoa(int(id))
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		return tasks, fmt.Errorf("ошибка перебора БД: %w", err)
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var (
		task   Task
		taskId int64
	)

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	err := DB.QueryRow(query, id).Scan(
		&taskId, &task.Date, &task.Title, &task.Comment, &task.Repeat,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("задача не найдена")
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка поучениЯ задачи: %w", err)
	}
	task.ID = strconv.Itoa(int(taskId))
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибкка обновления задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка оновления: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки удаления: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateDate(id, next string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки обновления даты: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
