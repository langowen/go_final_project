package tests

import (
	"github.com/langowen/go_final_project/internal/config"
	"path/filepath"
)

const pathEnv = "../.env"

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = false
var Token = ``

func init() {
	loadConfig()
}

func loadConfig() {
	cfg := config.MustLoad(pathEnv)
	Port = cfg.Port
	DBFile = filepath.Join("..", cfg.FileDb)
}
