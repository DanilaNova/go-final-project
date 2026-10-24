package api

import (
	"log"
	"net/http"

	"github.com/DanilaNova/go-final-project/pkg/db"
)

const taskGetLimit = 50

type TaskList struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandlerGet(response http.ResponseWriter, request *http.Request) {
	search := request.FormValue("search")

	tasks, err := db.GetTasks(taskGetLimit, search)
	if err != nil {
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
