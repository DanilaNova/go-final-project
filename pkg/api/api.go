package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJson(response http.ResponseWriter, code int, data any) error {
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response.WriteHeader(code)

	data_json, err := json.Marshal(data)
	if err == nil {
		_, err = response.Write(data_json)
	}
	return err
}

func writeError(response http.ResponseWriter, code int, err error) {
	err = writeJson(response, code, ErrorResponse{err.Error()})
	if err != nil {
		log.Println("ERROR: could not write error response: ", err)
	}
}

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/signin", sighinHandler)
}
