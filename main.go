package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var debug = os.Getenv("DEBUG") == "true"

func main() {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = ":8080"
	}

	setupHandlers()
	setupS3()

	router := mux.NewRouter()
	router.Use(middlewareLogger)

	router.Methods(http.MethodGet).Path("/readyz").HandlerFunc(readyz)
	router.Methods(http.MethodGet).PathPrefix("/assets").HandlerFunc(handleAssets)

	router.Methods(http.MethodGet).Path("/").HandlerFunc(handleCreate)
	router.Methods(http.MethodPost).Path("/").HandlerFunc(handleCreateForm)
	router.Methods(http.MethodGet).Path("/help").HandlerFunc(handleHelp)
	uploadRouter := router.PathPrefix("/{id}").Subrouter()

	uploadTemplateRouter := uploadRouter.Path("").Subrouter()
	s3Router := uploadRouter.PathPrefix("/s3/multipart").Subrouter()

	uploadTemplateRouter.Methods(http.MethodGet).Path("").HandlerFunc(handleUpload)

	s3Router.Methods(http.MethodPost).Path("").HandlerFunc(handleCreateMultipartUpload)
	s3Router.Methods(http.MethodGet).Path("/{uploadID}").HandlerFunc(handleGetUploadedParts)
	s3Router.Methods(http.MethodGet).Path("/{uploadID}/batch").HandlerFunc(handleBatchSignPartsUpload)
	s3Router.Methods(http.MethodGet).Path("/{uploadID}/{uploadPart}").HandlerFunc(handleSignPartUpload)
	s3Router.Methods(http.MethodPost).Path("/{uploadID}/complete").HandlerFunc(handleCompleteMultipartUpload)
	s3Router.Methods(http.MethodDelete).Path("/{uploadID}").HandlerFunc(handleAbortMultipartUpload)

	server := &http.Server{
		Handler:     router,
		Addr:        listen,
		ReadTimeout: 30 * time.Second,
	}
	log.Printf("listening on %s", listen)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
