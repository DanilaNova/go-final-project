package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/DanilaNova/go-final-project/pkg/db"
)

var ErrTitleIsEmpty = errors.New("task title is empty")
var ErrNoIdentifier = errors.New("identifier wasn't specified")
var ErrTaskNotFound = errors.New("task not found")

type TaskId struct {
	ID int64 `json:"id"`
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

func parseTask(request *http.Request) (*db.Task, error) {
	content, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println("ERROR: could not read request content: ", err)
		return nil, err
	}

	task := new(db.Task)
	err = json.Unmarshal(content, task)
	if err != nil {
		log.Println("ERROR: could not parse request content: ", err)
		return nil, err
	}

	if len(task.Title) == 0 {
		log.Println("ERROR: ", ErrTitleIsEmpty)
		return nil, ErrTitleIsEmpty
	}

	if err := checkDate(task); err != nil {
		log.Println("ERROR: error in task: ", err)
		return nil, err
	}

	return task, nil
}

func taskHandlerGet(response http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	if len(id) == 0 {
		log.Println("ERROR: ", ErrNoIdentifier)
		writeError(response, http.StatusInternalServerError, ErrNoIdentifier)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("WARNING: " + ErrTaskNotFound.Error() + "(id: \"" + id + "\")")
			writeError(response, http.StatusInternalServerError, ErrTaskNotFound)
			return
		}
		log.Println("ERROR: sql error: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	err = writeJson(response, task)
	if err != nil {
		log.Println("ERROR: could not write response: ", err)
		writeError(response, http.StatusInternalServerError, err)
	}
}

func taskHandlerPost(response http.ResponseWriter, request *http.Request) {
	task, err := parseTask(request)
	if err != nil {
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	id, err := db.AddTask(task)
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

func taskHandlerPut(response http.ResponseWriter, request *http.Request) {
	task, err := parseTask(request)
	if err != nil {
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	err = db.UpdateTask(task)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("WARNING: " + ErrTaskNotFound.Error() + "(id: \"" + task.ID + "\")")
			writeError(response, http.StatusInternalServerError, ErrTaskNotFound)
			return
		}
		log.Println("ERROR: could not update task: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	err = writeJson(response, struct{}{})
	if err != nil {
		log.Println("ERROR: could not write response: ", err)
	}
}

func taskHandlerDelete(response http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	if len(id) == 0 {
		log.Println("ERROR: ", ErrNoIdentifier)
		writeError(response, http.StatusInternalServerError, ErrNoIdentifier)
		return
	}

	_, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("WARNING: " + ErrTaskNotFound.Error() + "(id: \"" + id + "\")")
			writeError(response, http.StatusInternalServerError, ErrTaskNotFound)
			return
		}
		log.Println("ERROR: sql error: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		log.Println("ERROR: could not delete a task: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	response.WriteHeader(http.StatusOK)
	err = writeJson(response, struct{}{})
	if err != nil {
		log.Println("ERROR: could not write response: ", err)
		writeError(response, http.StatusInternalServerError, err)
	}
}

func taskHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		taskHandlerPost(response, request)
	case http.MethodGet:
		taskHandlerGet(response, request)
	case http.MethodPut:
		taskHandlerPut(response, request)
	case http.MethodDelete:
		taskHandlerDelete(response, request)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Header().Set("Allow", "GET, POST, PUT, DELETE")
	}
}
