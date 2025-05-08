package tests

import (
	"github.com/langowen/go_final_project/pkg/config"
	"path/filepath"
)

const pathEnv = "../.env"

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY3NzA1MDYsInB3ZF9oYXNoIjoiN2I3YmYxNmVhMjZjZTQ3MTVlNTZjYzQ2OWMwZjNjNzdkZTU3MWYxOTJjYTg4MTUyNDVhZWI2M2ZmMTc1YjQxOSJ9.CvOECcepXJZ-PyWhcFTGaPuRptjmNzy3GFHWLuUVZ0I"

func init() {
	loadConfig()
}

func loadConfig() {
	cfg := config.MustLoad(pathEnv)
	Port = cfg.Port
	DBFile = filepath.Join("..", cfg.FileDb)
}
