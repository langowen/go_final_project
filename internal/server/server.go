package server

import (
	"fmt"
	"github.com/langowen/go_final_project/internal/config"
	"net/http"
)

func NewServer(cfg *config.Config) error {

	serverPort := ":" + fmt.Sprintf("%d", cfg.Port)

	todo := http.NewServeMux()

	todo.Handle("/", http.FileServer(http.Dir(cfg.WebDir)))

	fmt.Printf("Сервер запущен, порт: %d", cfg.Port)

	return http.ListenAndServe(serverPort, todo)
}
