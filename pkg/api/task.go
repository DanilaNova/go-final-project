package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/DanilaNova/go-final-project/pkg/db"
)

var ErrTitleIsEmpty = errors.New("task title is empty")

type TaskId struct {
	ID int64 `json:"id"`
}

type TaskList struct {
	Tasks []*db.Task `json:"tasks"`
}

func checkDate(task *db.Task) error {
	if len(task.Date) == 0 {
		task.Date = time.Now().Format(db.DateFormat)
	}

	t, err := time.Parse(db.DateFormat, task.Date)
	if err != nil {
		return err
	}

	now := time.Now()
	var next string
	if len(task.Repeat) != 0 {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if dateIsAfter(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(db.DateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

func taskHandlerPost(response http.ResponseWriter, request *http.Request) {
	content, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println("ERROR: could not read request content: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	var task db.Task
	err = json.Unmarshal(content, &task)
	if err != nil {
		log.Println("ERROR: could not parse request content: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	if len(task.Title) == 0 {
		log.Println("ERROR: ", ErrTitleIsEmpty)
		writeError(response, http.StatusInternalServerError, ErrTitleIsEmpty)
		return
	}

	if err := checkDate(&task); err != nil {
		log.Println("ERROR: error in task: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Println("ERROR: could not add task: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	err = writeJson(response, TaskId{id})
	if err != nil {
		log.Println("ERROR: could not write response: ", err)
		writeError(response, http.StatusInternalServerError, err)
	}
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

func taskHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		taskHandlerPost(response, request)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Header().Set("Allow", "POST")
	}
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
