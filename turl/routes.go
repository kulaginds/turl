package turl

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "turl/turl/models"
)

type UrlResponse struct {
	Url string `json:"urls"`
}

func JsonResponse(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(r)
}

func HelloRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Tiny urls - urls shorter")
}

func ShortRoute(w http.ResponseWriter, r *http.Request) {
	url := ShortUrl{}

	json.NewDecoder(r.Body).Decode(&url)

	if 0 == len(url.Url) {
		http.NotFound(w, r)
		return
	}

	JsonResponse(w, &UrlResponse{Url:url.Url})
}

func LongRoute(w http.ResponseWriter, r *http.Request) {
	url := LongUrl{}

	json.NewDecoder(r.Body).Decode(&url)

	if 0 == len(url.Url) {
		http.NotFound(w, r)
		return
	}

	JsonResponse(w, &UrlResponse{Url:url.Url})
}
