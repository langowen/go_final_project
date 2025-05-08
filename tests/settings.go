package tests

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/langowen/go_final_project/pkg/auth"
	"github.com/langowen/go_final_project/pkg/config"
	"path/filepath"
	"time"
)

const pathEnv = "../.env"

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = ""

func init() {
	loadConfig()
}

func init() {
	cfg := loadConfig()
	generateTestToken(cfg)
}

func loadConfig() *config.Config {
	cfg := config.MustLoad(pathEnv)
	Port = cfg.Port
	DBFile = filepath.Join("..", cfg.FileDb)
	return cfg
}

func generateTestToken(cfg *config.Config) {
	if cfg.Token == "" {
		Token = ""
		return
	}

	claims := jwt.MapClaims{
		"exp":      time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC).Unix(), // Фиксированная дата
		"pwd_hash": auth.HashPassword(cfg.Token),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Token))
	if err != nil {
		panic("Failed to generate test token: " + err.Error())
	}

	Token = tokenString
}
