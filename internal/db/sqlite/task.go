package sqlite

import (
	"fmt"
	"github.com/langowen/go_final_project/internal/db"
	"strconv"
	"time"
)

func (s *Storage) AddTask(task *db.Task) (string, error) {
	// Проверка обязательных полей
	if task.Title == "" {
		return "", fmt.Errorf("title is required")
	}

	// Валидация даты
	if task.Date != "" {
		if _, err := time.Parse("20060102", task.Date); err != nil {
			return "", fmt.Errorf("invalid date format")
		}
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", fmt.Errorf("database error: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return strconv.FormatInt(id, 10), nil
}
