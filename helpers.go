package main

import (
	"errors"
	"net/http"
)

var errNotFound = errors.New("not found")
var errBadRequest = errors.New("bad request")
var errInternalServerError = errors.New("internal server error")
var errUnauthorized = errors.New("unauthorized")
var errForbidden = errors.New("forbidden")
var errConflict = errors.New("conflict")

func errorResponseStatus(w http.ResponseWriter, req *http.Request, err error) {
	errorStatus := http.StatusInternalServerError

	if errors.Is(err, errNotFound) {
		errorStatus = http.StatusNotFound
	} else if errors.Is(err, errBadRequest) {
		errorStatus = http.StatusBadRequest
	} else if errors.Is(err, errInternalServerError) {
		errorStatus = http.StatusInternalServerError
	} else if errors.Is(err, errUnauthorized) {
		errorStatus = http.StatusUnauthorized
	} else if errors.Is(err, errForbidden) {
		errorStatus = http.StatusForbidden
	} else if errors.Is(err, errConflict) {
		errorStatus = http.StatusConflict
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

// responseToError converts a HTTP status code to an error
func responseToError(resp *http.Response) error {
	if resp.StatusCode == http.StatusNotFound {
		return errNotFound
	} else if resp.StatusCode == http.StatusBadRequest {
		return errBadRequest
	} else if resp.StatusCode == http.StatusInternalServerError {
		return errInternalServerError
	} else if resp.StatusCode == http.StatusUnauthorized {
		return errUnauthorized
	} else if resp.StatusCode == http.StatusForbidden {
		return errForbidden
	} else if resp.StatusCode == http.StatusConflict {
		return errConflict
	}
	return nil
}
