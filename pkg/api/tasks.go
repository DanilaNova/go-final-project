package api

import (
	"log"
	"net/http"

	"github.com/DanilaNova/go-final-project/pkg/db"
)

type TaskList struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandlerGet(response http.ResponseWriter, request *http.Request) {
	search := request.FormValue("search")

	tasks, err := db.GetTasks(50, search)
	if err != nil {
		log.Println("ERROR: cound not get tasks: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	writeJson(response, TaskList{tasks})
}

func tasksHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		tasksHandlerGet(response, request)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Header().Set("Allow", "GET")
	}
}
