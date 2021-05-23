package main

import (
	"io/fs"
	"log"
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

	setupHandlers()
	setupS3()

	router := mux.NewRouter()
	router.Use(middlewareLogger)

	router.Methods(http.MethodGet).Path("/readyz").HandlerFunc(readyz)
	router.Methods(http.MethodGet).PathPrefix("/assets").Handler(http.FileServer(http.FS(assetsWeb)))
	router.Methods(http.MethodGet).Path("/").HandlerFunc(handleCreate)
	uploadRouter := router.PathPrefix("/{id}").Subrouter()

	uploadTemplateRouter := uploadRouter.Path("").Subrouter()
	s3Router := uploadRouter.PathPrefix("/s3/multipart").Subrouter()

	uploadTemplateRouter.Methods(http.MethodGet).Path("").HandlerFunc(handleUpload)

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
	log.Printf("listeining on %s", listen)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
