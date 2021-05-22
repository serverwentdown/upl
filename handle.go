package main

import (
	"net/http"
	"os"
)

func getUploadedParts(w http.ResponseWriter, req *http.Request) {
}

func signPartUpload(w http.ResponseWriter, req *http.Request) {
	method := http.MethodGet
	url := "https://minio1.makerforce.io/test"

	unsignedReq, err := http.NewRequest(method, url, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	cred := credential{
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		Region:    os.Getenv("MINIO_REGION_NAME"),
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),
	}

	signedReq := preSign(unsignedReq, cred)

	w.Write([]byte(signedReq.URL.String()))
}

func completeMultipartUpload(w http.ResponseWriter, req *http.Request) {
}

func abortMultipartUpload(w http.ResponseWriter, req *http.Request) {
}

var globalStore store

func setupHandlers() {
	var err error
	globalStore, err = newRedisStore(os.Getenv("REDIS_CONNECTION"))
	if err != nil {
		panic(err)
	}
}
