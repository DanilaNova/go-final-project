package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/DanilaNova/go-final-project/pkg/db"
)

const TASK_GET_LIMIT = 50

type TaskList struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandlerGet(response http.ResponseWriter, request *http.Request) {
	search := request.FormValue("search")

	tasks, err := db.GetTasks(TASK_GET_LIMIT, search)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("WARNING: got no tasks")
			writeError(response, http.StatusNotFound, err)
			return
		}
		log.Println("ERROR: cound not get tasks: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	err = writeJson(response, http.StatusOK, TaskList{tasks})
	if err != nil {
		log.Println("ERROR: could not write response: ", err)
	}
}

func tasksHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		tasksHandlerGet(response, request)
	default:
		response.Header().Set("Allow", "GET")
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}
