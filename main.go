package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"turl/turl"
)

func main() {
	fmt.Println("Tiny url")

	r := mux.NewRouter()

	r.HandleFunc("/", turl.HelloRoute)
	r.HandleFunc("/short", turl.ShortRoute).Methods("POST")
	r.HandleFunc("/long", turl.LongRoute).Methods("POST")
	r.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
