package main

import (
	"context"
	"fmt"
	"github.com/langowen/go_final_project/internal/config"
	"github.com/langowen/go_final_project/internal/db/sqlite"
	"github.com/langowen/go_final_project/internal/server"
	"time"
)

const (
	pathEnv = ".env"
)

func main() {
	cfg := config.MustLoad(pathEnv)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sqlite3, err := sqlite.Init(ctx, cfg.FileDb)
	if err != nil {
		fmt.Println(err)
	}

	err = server.NewServer(cfg, sqlite3)
	if err != nil {
		fmt.Printf("не удалось запустить сервер: %s", err)
	}
}
