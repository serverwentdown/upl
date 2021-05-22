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
	multipartRouter := router.PathPrefix("/s3/multipart").Subrouter()
	router.PathPrefix("/").Handler(http.FileServer(http.FS(assetsWeb)))

	multipartRouter.HandleFunc("", handleCreateMultipartUpload).Methods(http.MethodPost)
	multipartRouter.HandleFunc("/{id}", handleGetUploadedParts).Methods(http.MethodGet)
	multipartRouter.HandleFunc("/{id}/{part}", handleSignPartUpload).Methods(http.MethodGet)
	multipartRouter.HandleFunc("/{id}/complete", handleCompleteMultipartUpload).Methods(http.MethodPost)
	multipartRouter.HandleFunc("/{id}", handleAbortMultipartUpload).Methods(http.MethodDelete)

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
