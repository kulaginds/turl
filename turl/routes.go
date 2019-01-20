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

type UrlError struct {
	Code        uint    `json:"code"`
	Description string `json:"description"`
}

func JsonResponse(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(r)
}

func JsonClientError(w http.ResponseWriter, errorCode uint) {
	w.WriteHeader(http.StatusBadRequest)
	JsonResponse(w, &UrlError{Code:errorCode, Description:ErrMsg[errorCode]})
}

func JsonServerError(w http.ResponseWriter, errorCode uint) {
	w.WriteHeader(http.StatusInternalServerError)
	JsonResponse(w, &UrlError{Code:errorCode, Description:ErrMsg[errorCode]})
}

func HelloRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Tiny urls - urls shorter")
}

func ShortRoute(w http.ResponseWriter, r *http.Request) {
	var longUrl *LongUrl
	url := ShortUrl{}

	json.NewDecoder(r.Body).Decode(&url)

	code, ok := url.Validate()

	if !ok {
		JsonClientError(w, code)
		return
	}

	longUrl, ok = url.Short()

	if !ok {
		JsonServerError(w, ErrShortLinkFail)
		return
	}

	JsonResponse(w, &UrlResponse{Url:longUrl.Url})
}

func LongRoute(w http.ResponseWriter, r *http.Request) {
	url := LongUrl{}

	json.NewDecoder(r.Body).Decode(&url)

	code, ok := url.Validate()

	if !ok {
		JsonClientError(w, code)
		return
	}

	JsonResponse(w, &UrlResponse{Url:url.Url})
}
