package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port   int    `env:"TODO_PORT" env-required:"true" `
	WebDir string `env:"TODO_WEB_DIR" env-default:"./web/"`
	FileDb string `env:"TODO_DBFILE" env-required:"true"`
}

func MustLoad(pathEnv string) *Config {
	if err := godotenv.Load(pathEnv); err != nil {
		fmt.Println("failed to load .env file: " + err.Error())
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read env vars: " + err.Error())
	}

	return &cfg
}
