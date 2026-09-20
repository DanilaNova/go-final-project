package server

import (
	"net/http"
	"strconv"

	"github.com/DanilaNova/go-final-project/pkg/api"
)

func Serve(webDir string, port int) error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init()

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
