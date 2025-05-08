package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func Init(ctx context.Context, path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть базу данных: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("не удалось установить соединение с базой данных: %w", err)
	}

	storage := &Storage{db: db}

	if err = storage.new(ctx); err != nil {
		return nil, fmt.Errorf("не удалось создать таблицы в базе данных: %w", err)
	}

	return storage, nil
}

func (s *Storage) new(ctx context.Context) error {
	q := `
    CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date DATE,
        title VARCHAR,
        comment TEXT,
        repeat VARCHAR
    );`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);",
	}

	for _, index := range indexes {
		_, err = s.db.ExecContext(ctx, index)
		if err != nil {
			return err
		}
	}

	return nil
}
