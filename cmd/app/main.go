package main

import (
	"context"
	"github.com/langowen/go_final_project/pkg/config"
	"github.com/langowen/go_final_project/pkg/db/sqlite"
	"github.com/langowen/go_final_project/pkg/server"
	"log"
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
		log.Fatal(err)
	}

	err = server.NewServer(cfg, sqlite3)
	if err != nil {
		log.Fatalf("не удалось запустить сервер: %s", err)
	}
}
