package main

import (
	"errors"
	"log"
	"net/http"
)

var errNotFound = errors.New("not found")
var errBadRequest = errors.New("bad request")
var errInternalServerError = errors.New("internal server error")

func errorResponseStatus(w http.ResponseWriter, req *http.Request, err error) {
	errorMessage := err.Error()
	errorStatus := http.StatusInternalServerError

	if errors.Is(err, errNotFound) {
		errorStatus = http.StatusNotFound
	} else if errors.Is(err, errBadRequest) {
		errorStatus = http.StatusBadRequest
	} else if errors.Is(err, errInternalServerError) {
		errorStatus = http.StatusInternalServerError
	}

	log.Printf("%s %s: %s", req.Method, req.RequestURI, errorMessage)
	w.WriteHeader(errorStatus)
}

func errorResponse(w http.ResponseWriter, req *http.Request, err error) {
	errorResponseStatus(w, req, err)
	w.Write([]byte(err.Error()))
}
