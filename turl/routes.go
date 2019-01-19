package turl

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UrlResponse struct {
	Url string `json:"url"`
}

func jsonResponse(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(r)
}

func HelloRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Tiny url - url shorter")
}

func ShortRoute(w http.ResponseWriter, r *http.Request) {
	var url Url

	json.NewDecoder(r.Body).Decode(&url)

	if 0 == len(url.Url) {
		http.NotFound(w, r)
		return
	}

	jsonResponse(w, &UrlResponse{Url:url.Url})
}

func LongRoute(w http.ResponseWriter, r *http.Request) {
	var url Url

	json.NewDecoder(r.Body).Decode(&url)

	if 0 == len(url.Url) {
		http.NotFound(w, r)
		return
	}

	jsonResponse(w, &UrlResponse{Url:url.Url})
}
