package turl

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"turl/turl/models"
)

var config *models.Config
var router *mux.Router

var initComplete bool

func Initialize() {
	fmt.Println("Tiny urls")

	config = models.NewConfig()
	router = mux.NewRouter()

	if !models.Initialize(config) {
		return
	}

	router.HandleFunc("/", HelloRoute)
	router.HandleFunc("/short", ShortRoute).Methods("POST")
	router.HandleFunc("/long", LongRoute).Methods("POST")
	router.Use(loggingMiddleware)

	initComplete = true
}

func Run() {
	if !initComplete {
		return
	}

	fmt.Println("Service url:", config.ServiceUrl())
	fmt.Println("Service started on port:", config.ServicePort())
	log.Fatal(http.ListenAndServe(config.ServicePort(), router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
