package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/langowen/go_final_project/pkg/db"
	"strconv"
	"time"
)

func (s *Storage) AddTask(task *db.Task) (string, error) {
	if task.Title == "" {
		return "", fmt.Errorf("title is required")
	}

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

func (s *Storage) Tasks(limit int) ([]*db.Task, error) {
	today := time.Now().Format("20060102")
	const query = `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE date >= ?
		ORDER BY date, id
		LIMIT ?
	`

	rows, err := s.db.Query(query, today, limit)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func (s *Storage) SearchTasks(search string, limit int) ([]*db.Task, error) {
	if date, err := time.Parse("02.01.2006", search); err == nil {
		return s.searchByDate(date.Format("20060102"), limit)
	}

	return s.searchByText(search, limit)
}

func (s *Storage) Task(id string) (*db.Task, error) {
	const query = `
        SELECT id, date, title, comment, repeat 
        FROM scheduler 
        WHERE id = ?
    `

	var task db.Task

	err := s.db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.ErrTaskNotFound
		}

		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}

	return &task, nil
}

func (s *Storage) UpdateTask(task *db.Task) error {
	const query = `
		UPDATE scheduler 
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func (s *Storage) Delete(id string) error {
	const query = `
		DELETE FROM scheduler 
		WHERE id = ?
	`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}
	return nil
}

func (s *Storage) searchByDate(date string, limit int) ([]*db.Task, error) {
	const query = `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE date = ?
		ORDER BY date, id
		LIMIT ?
	`

	rows, err := s.db.Query(query, date, limit)
	if err != nil {
		return nil, fmt.Errorf("query tasks by date: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func (s *Storage) searchByText(search string, limit int) ([]*db.Task, error) {
	const query = `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE title LIKE '%' || ? || '%' OR comment LIKE '%' || ? || '%'
		AND date >= ?
		ORDER BY date, id
		LIMIT ?
	`

	today := time.Now().Format("20060102")
	rows, err := s.db.Query(query, search, search, today, limit)
	if err != nil {
		return nil, fmt.Errorf("query tasks by text: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]*db.Task, error) {
	var tasks []*db.Task
	for rows.Next() {
		var t db.Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tasks, nil
}
