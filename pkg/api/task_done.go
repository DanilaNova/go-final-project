package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/DanilaNova/go-final-project/pkg/db"
)

func taskDoneHandlerPost(response http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	if len(id) == 0 {
		log.Println("ERROR: ", ErrNoIdentifier)
		writeError(response, http.StatusBadRequest, ErrNoIdentifier)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("WARNING: " + ErrTaskNotFound.Error() + "(id: \"" + id + "\")")
			writeError(response, http.StatusNotFound, ErrTaskNotFound)
			return
		}
		log.Println("ERROR: sql error: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	if len(task.Repeat) == 0 {
		err = db.DeleteTask(id)
		if err != nil {
			log.Println("ERROR: could not delete a task: ", err)
			writeError(response, http.StatusInternalServerError, err)
			return
		}

		err = writeJson(response, http.StatusOK, struct{}{})
		if err != nil {
			log.Println("ERROR: could not write response: ", err)
			writeError(response, http.StatusInternalServerError, err)
		}
		return
	}

	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		log.Println("ERROR: cannot solve next date: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	err = db.UpdateDate(next, id)
	if err != nil {
		log.Println("ERROR: cannot update date: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	err = writeJson(response, http.StatusOK, struct{}{})
	if err != nil {
		log.Println("ERROR: could not write response: ", err)
		writeError(response, http.StatusInternalServerError, err)
	}
}

func taskDoneHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		taskDoneHandlerPost(response, request)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Header().Set("Allow", "POST")
	}
}
