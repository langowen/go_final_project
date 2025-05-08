package server

import (
	"fmt"
	"github.com/langowen/go_final_project/internal/api"
	"github.com/langowen/go_final_project/internal/config"
	"github.com/langowen/go_final_project/internal/db"
	"net/http"
)

func NewServer(cfg *config.Config, storage db.Storage) error {
	serverPort := fmt.Sprintf(":%d", cfg.Port)

	todo := http.NewServeMux()
	todo.Handle("/", http.FileServer(http.Dir(cfg.WebDir)))

	api.Init(todo, storage)

	fmt.Printf("Сервер запущен, порт: %d\n", cfg.Port)

	if err := http.ListenAndServe(serverPort, todo); err != nil {
		return fmt.Errorf("ошибка запуска сервера: %w", err)
	}

	return nil
}
