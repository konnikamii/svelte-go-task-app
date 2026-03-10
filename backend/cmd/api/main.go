package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/handlers"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetReportCaller(true)
	var router *chi.Mux = chi.NewRouter()
	handlers.Handler(router)

	fmt.Println("Starting server....")

	err := http.ListenAndServe("localhost:8000", router)
	if err != nil {
		log.Error(err)
	}

}
