package main

import (
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = ":8080"
	}

	assetsWeb, err := fs.Sub(assets, "web")
	if err != nil {
		panic(err)
	}

	setupS3()
	setupHandlers()

	router := mux.NewRouter()
	uploadRouter := router.PathPrefix("/{id}").Subrouter()
	router.PathPrefix("/assets").Handler(http.FileServer(http.FS(assetsWeb)))

	uploadRouter.Path("").HandlerFunc(handleUpload)
	s3Router := uploadRouter.PathPrefix("/s3/multipart").Subrouter()

	s3Router.Methods(http.MethodPost).Path("").HandlerFunc(handleCreateMultipartUpload)
	s3Router.Methods(http.MethodGet).Path("/{uploadID}").HandlerFunc(handleGetUploadedParts)
	s3Router.Methods(http.MethodGet).Path("/{uploadID}/{uploadPart}").HandlerFunc(handleSignPartUpload)
	s3Router.Methods(http.MethodPost).Path("/{uploadID}/complete").HandlerFunc(handleCompleteMultipartUpload)
	s3Router.Methods(http.MethodDelete).Path("/{uploadID}").HandlerFunc(handleAbortMultipartUpload)

	server := &http.Server{
		Handler:      router,
		Addr:         listen,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
