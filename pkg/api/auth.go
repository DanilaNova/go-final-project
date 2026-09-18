package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var ErrAuthentificationRequired = errors.New("authentification required")

type PassRequest struct {
	Password string `json:"password"`
}

type PassResponse struct {
	Token string `json:"token"`
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtStr string // JWT-токен из куки
			var valid bool

			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtStr = cookie.Value
				jwtToken, err := jwt.Parse(jwtStr, func(t *jwt.Token) (any, error) {
					return []byte(pass), nil
				})
				if err != nil {
					log.Println("ERROR: Failed to parse token: ", err)
					writeError(w, http.StatusInternalServerError, err)
					return
				}
				valid = jwtToken.Valid
			}

			if !valid {
				// возвращаем ошибку авторизации 401
				writeError(w, http.StatusUnauthorized, ErrAuthentificationRequired)
				return
			}
		}
		next(w, r)
	})
}

func sighinHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Header().Set("Allow", "POST")
		return
	}

	pass := os.Getenv("TODO_PASSWORD")

	data, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println("ERROR: could not read request content: ", err)
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	passRequest := new(PassRequest)
	err = json.Unmarshal(data, passRequest)
	if err != nil {
		log.Println("ERROR: could not parse pass request: ", err)
		writeError(response, http.StatusInternalServerError, err)
	}

	if passRequest.Password != string(pass) {
		log.Println("WARNING: authentification wasn't passed")
		writeError(response, http.StatusUnauthorized, ErrAuthentificationRequired)
		return
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	signedToken, err := jwtToken.SignedString([]byte(pass))
	if err != nil {
		log.Println("ERROR: failed to sign a token: ", err)
		writeError(response, http.StatusUnauthorized, ErrAuthentificationRequired)
		return
	}

	response.WriteHeader(http.StatusOK)
	err = writeJson(response, PassResponse{signedToken})
	if err != nil {
		log.Println("ERROR: could not write response: ", err)
		writeError(response, http.StatusInternalServerError, err)
	}
}
