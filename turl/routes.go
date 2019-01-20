package turl

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "turl/turl/models"
)

type Error struct {
	Code        uint    `json:"code"`
	Description string `json:"description"`
}

func JsonResponse(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(r)
}

func JsonError(w http.ResponseWriter, errorCode uint) {
	if errorCode < 10 {
		jsonClientError(w, errorCode)
	} else {
		jsonServerError(w, errorCode)
	}
}

func jsonClientError(w http.ResponseWriter, errorCode uint) {
	w.WriteHeader(http.StatusBadRequest)
	JsonResponse(w, &Error{Code:errorCode, Description:ErrMsg[errorCode]})
}

func jsonServerError(w http.ResponseWriter, errorCode uint) {
	w.WriteHeader(http.StatusInternalServerError)
	JsonResponse(w, &Error{Code:errorCode, Description:ErrMsg[errorCode]})
}

func HelloRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Tiny urls - urls shorter")
}

func ShortRoute(w http.ResponseWriter, r *http.Request) {
	longUrl := LongUrl{}

	json.NewDecoder(r.Body).Decode(&longUrl)

	shortUrl, code, ok := longUrl.Short()

	if !ok {
		JsonError(w, code)
		return
	}

	JsonResponse(w, shortUrl)
}

func LongRoute(w http.ResponseWriter, r *http.Request) {
	shortUrl := ShortUrl{}

	json.NewDecoder(r.Body).Decode(&shortUrl)

	longUrl, code, ok := shortUrl.Long()

	if !ok {
		JsonError(w, code)
		return
	}

	JsonResponse(w, longUrl)
}
