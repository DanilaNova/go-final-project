package server

import (
	"net/http"
	"strconv"
)

func Serve(webDir string, port int) error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	return http.ListenAndServe("localhost:"+strconv.Itoa(port), nil)
}
