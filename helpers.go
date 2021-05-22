package main

import (
	"errors"
	"log"
	"net/http"
)

var errNotFound = errors.New("not found")
var errBadRequest = errors.New("bad request")

func errorResponse(w http.ResponseWriter, req *http.Request, err error) {
	errorMessage := err.Error()
	errorStatus := http.StatusInternalServerError

	if errors.Is(err, errNotFound) {
		errorStatus = http.StatusNotFound
	} else if errors.Is(err, errBadRequest) {
		errorStatus = http.StatusBadRequest
	}

	log.Printf("%s %s: %s", req.Method, req.URL.Path, errorMessage)
	w.WriteHeader(errorStatus)
}
