package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func formatKey(prefix, filename string) string {
	for strings.Contains(prefix, "{random}") {
		random := gonanoid.MustGenerate(idAlphabet, 16)
		prefix = strings.Replace(prefix, "{random}", random, 1)
	}
	return prefix + filename
}

/* createMultipartUpload */

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
	vars := mux.Vars(req)
	cred, err := getCredential(vars["id"])
	if err != nil {
		errorResponse(w, req, err)
		return
	}

	r := createMultipartUploadReq{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&r); err != nil {
		errorResponse(w, req, fmt.Errorf("%w: %s", errBadRequest, err))
		return
	}

	if err := r.validate(); err != nil {
		errorResponse(w, req, err)
		return
	}

	// Derive the object key
	key := formatKey(cred.Prefix, r.Filename)

	result, err := initiateMultipartUpload(key, cred)
	if err != nil {
		errorResponse(w, req, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(createMultipartUploadRes{
		Key:      key,
		UploadID: result.UploadID,
	})
}

/* getUploadedParts */

type getUploadedPartsRes []part

func handleGetUploadedParts(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cred, err := getCredential(vars["id"])
	if err != nil {
		errorResponse(w, req, err)
		return
	}

	uploadID := vars["uploadID"]
	key := req.URL.Query().Get("key")

	if uploadID == "" || key == "" {
		errorResponse(w, req, fmt.Errorf("%w", errBadRequest))
		return
	}

	parts := make(getUploadedPartsRes, 0, 0)
	var nextPartNumberMarker uint32
	for {
		page, err := listParts(key, uploadID, cred, nextPartNumberMarker)
		if err != nil {
			errorResponse(w, req, err)
			return
		}

		parts = append(parts, page.Parts...)
		nextPartNumberMarker = page.NextPartNumberMarker

		if !page.IsTruncated {
			break
		}
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(getUploadedPartsRes(parts))
}

/* signPartUpload */

type signPartUploadRes struct {
	URL string `json:"url"`
}

func handleSignPartUpload(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cred, err := getCredential(vars["id"])
	if err != nil {
		errorResponse(w, req, err)
		return
	}

	uploadID := vars["uploadID"]
	key := req.URL.Query().Get("key")
	partNumber, err := strconv.ParseUint(vars["uploadPart"], 10, 16)

	if uploadID == "" || key == "" {
		errorResponse(w, req, fmt.Errorf("%w", errBadRequest))
		return
	}
	if partNumber < 1 || partNumber > 10000 || err != nil {
		errorResponse(w, req, fmt.Errorf("%w: invalid part number", errBadRequest))
		return
	}

	params := make(url.Values)
	params.Add("partNumber", strconv.FormatUint(partNumber, 10))
	params.Add("uploadId", uploadID)
	unsignedReq, err := http.NewRequest(http.MethodPut, cred.Endpoint+"/"+key+"?"+params.Encode(), nil)
	if err != nil {
		errorResponse(w, req, fmt.Errorf("%w: %s", errInternalServerError, err))
		return
	}

	signedReq := preSign(unsignedReq, cred)

	encoder := json.NewEncoder(w)
	encoder.Encode(signPartUploadRes{
		URL: signedReq.URL.String(),
	})
}

/* completeMultipartUpload */

type completeMultipartUploadReq struct {
	Parts []completePart `json:"parts"`
}

func (r completeMultipartUploadReq) validate() error {
	for _, part := range r.Parts {
		if err := part.validate(); err != nil {
			return err
		}
	}
	return nil
}

func handleCompleteMultipartUpload(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cred, err := getCredential(vars["id"])
	if err != nil {
		errorResponse(w, req, err)
		return
	}

	uploadID := vars["uploadID"]
	key := req.URL.Query().Get("key")

	if uploadID == "" || key == "" {
		errorResponse(w, req, fmt.Errorf("%w", errBadRequest))
		return
	}

	r := completeMultipartUploadReq{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&r); err != nil {
		errorResponse(w, req, fmt.Errorf("%w: %s", errBadRequest, err))
		return
	}

	if err := r.validate(); err != nil {
		errorResponse(w, req, err)
		return
	}

	result, err := completeMultipartUpload(key, uploadID, r.Parts, cred)
	if err != nil {
		errorResponse(w, req, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

/* abortMultipartUpload */

func handleAbortMultipartUpload(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cred, err := getCredential(vars["id"])
	if err != nil {
		errorResponse(w, req, err)
		return
	}

	uploadID := vars["uploadID"]
	key := req.URL.Query().Get("key")

	if uploadID == "" || key == "" {
		errorResponse(w, req, fmt.Errorf("%w", errBadRequest))
		return
	}

	err = abortMultipartUpload(key, uploadID, cred)
	if err != nil {
		errorResponse(w, req, err)
		return
	}
}
