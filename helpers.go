package main

import (
	"errors"
	"net/http"
)

var errNotFound = errors.New("not found")
var errBadRequest = errors.New("bad request")
var errInternalServerError = errors.New("internal server error")

func errorResponseStatus(w http.ResponseWriter, req *http.Request, err error) {
	errorStatus := http.StatusInternalServerError

	if errors.Is(err, errNotFound) {
		errorStatus = http.StatusNotFound
	} else if errors.Is(err, errBadRequest) {
		errorStatus = http.StatusBadRequest
	} else if errors.Is(err, errInternalServerError) {
		errorStatus = http.StatusInternalServerError
	}

	w.WriteHeader(errorStatus)
}

// errorResponse prints the error message in the response body.
//
// Do not use this function when the error message might contain sensitive
// information.
func errorResponse(w http.ResponseWriter, req *http.Request, err error) {
	errorResponseStatus(w, req, err)
	w.Write([]byte(err.Error()))
}
