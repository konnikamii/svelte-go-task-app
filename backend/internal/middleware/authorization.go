package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/konnikamii/svelte-go-task-app/backend/api"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/tools"
	log "github.com/sirupsen/logrus"
)

var UnauthorizedError = errors.New("Invalid username or token...")

type Book struct {
	Id   int32  `json:"id,omitempty" bson:"id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body Book
		var err error
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return
		}
		println("body", body)
		var username string = r.URL.Query().Get("username")
		var token string = r.Header.Get("Authorization")

		if username == "" || token == "" {
			log.Error(UnauthorizedError)
			api.RequestErrorHandler(w, UnauthorizedError)
			return
		}

		var database *tools.DatabaseInterface
		database, err = tools.NewDatabase()
		if err != nil {
			api.InternalErrorHandler(w)
			return
		}

		var loginDetails *tools.LoginDetails
		loginDetails = (*database).GetUserLoginDetails(username)

		if loginDetails == nil || (token != (*loginDetails).AuthToken) {
			log.Error(UnauthorizedError)
			api.RequestErrorHandler(w, UnauthorizedError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
