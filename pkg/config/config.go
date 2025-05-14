package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port   int    `env:"TODO_PORT" env-default:"8088" `
	WebDir string `env:"TODO_WEB_DIR" env-default:"./web/"`
	FileDb string `env:"TODO_DBFILE" env-default:"scheduler.db"`
	Token  string `env:"TODO_PASSWORD"`
}

func MustLoad(pathEnv string) *Config {
	godotenv.Load(pathEnv)

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		fmt.Println("не удалось прочитать env переменные: " + err.Error())
	}

	return &cfg
}
