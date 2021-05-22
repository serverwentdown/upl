package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type createMultipartUploadReq struct {
	Filename string                           `json:"filename"`
	Type     string                           `json:"type"`
	Metadata createMultipartUploadReqMetadata `json:"metadata"`
}

func (r createMultipartUploadReq) validate() error {
	if r.Filename == "" {
		return errors.New("invalid filename")
	} else if r.Type == "" {
		return errors.New("invalid content type")
	}
	return nil
}

type createMultipartUploadReqMetadata struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type createMultipartUploadRes struct {
	Key      string `json:"key"`
	UploadID string `json:"uploadId"`
}

func handleCreateMultipartUpload(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	r := createMultipartUploadReq{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&r); err != nil {
		errorResponse(w, req, fmt.Errorf("%w: %s", errBadRequest, err))
		return
	}

	if err := r.validate(); err != nil {
		errorResponse(w, req, fmt.Errorf("%w: %s", errBadRequest, err))
		return
	}

	// Derive the object key
	// TODO: configurable
	key := fmt.Sprintf("uploads/%s", r.Filename)

	cred := credential{
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		Region:    os.Getenv("MINIO_REGION_NAME"),
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),
	}
	uploadID, err := createMultipartUpload(key, cred)
	if err != nil {
		errorResponse(w, req, fmt.Errorf("%w: %s", errBadRequest, err))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(createMultipartUploadRes{
		Key:      key,
		UploadID: uploadID,
	})
}
