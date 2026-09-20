package main

import (
	"log"
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
		eport, err := strconv.Atoi(envPort)
		if err != nil {
			log.Println("ERROR: could not parse port setting (\"" + envPort + "\"")
		} else {
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

	err := db.Init(dbFile)
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println("ERROR: could not close database connection: ", err)
		}
	}()
	if err != nil {
		panic(err)
	}

	log.Println("INFO: starting server")
	if err := server.Serve(webDir, port); err != nil {
		panic(err)
	}
}
