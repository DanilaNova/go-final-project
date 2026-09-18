package main

import (
	"os"
	"strconv"

	"github.com/DanilaNova/go-final-project/pkg/db"
	"github.com/DanilaNova/go-final-project/pkg/server"
)

// Defaults
var port = 7540
var webDir = "./web"
var dbFile = "./scheduler.db"

func main() {
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if eport, err := strconv.Atoi(envPort); err != nil {
			port = eport
		}
	}

	envDbFile := os.Getenv("TODO_DBFILE")
	if len(envDbFile) > 0 {
		dbFile = envDbFile
	}

	envWebDir := os.Getenv("TODO_WEBDIR")
	if len(envWebDir) > 0 {
		webDir = envWebDir
	}

	var db db.Database
	err := db.Init(dbFile)
	if err != nil {
		panic(err)
	}

	if err := server.Serve(webDir, port); err != nil {
		panic(err)
	}
}
