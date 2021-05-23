package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var globalStore store

func setupHandlers() {
	var err error
	globalStore, err = newRedisStore(os.Getenv("REDIS_CONNECTION"))
	if err != nil {
		panic(err)
	}

	if debug {
		assetsServer = http.FileServer(http.FS(os.DirFS("web")))
	}
}

/* assets */

var assetsServer = http.FileServer(http.FS(assetsWeb))

func handleAssets(w http.ResponseWriter, req *http.Request) {
	assetsServer.ServeHTTP(w, req)
}

/* templates */

var tmpl = template.Must(template.ParseFS(assets, "web/*.tmpl"))

func executeTemplate(w io.Writer, name string, data interface{}) error {
	if debug {
		tmpl = template.Must(template.ParseGlob("web/*.tmpl"))
	}
	return tmpl.ExecuteTemplate(w, name, nil)
}

/* credentials */

var idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

func getCredential(id string) (credential, error) {
	cred := credential{}

	b, err := globalStore.get(id)
	if err != nil {
		return cred, err
	}

	err = json.Unmarshal(b, &cred)
	if err != nil {
		return cred, err
	}

	return cred, nil
}

func setCredential(id string, cred credential, expire time.Duration) error {
	b, err := json.Marshal(cred)
	if err != nil {
		return err
	}

	err = globalStore.put(id, b, expire)
	if err != nil {
		return err
	}

	return nil
}

/* upload template */

func handleUpload(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	_, err := getCredential(vars["id"])
	if errors.Is(err, errNotFound) {
		errorResponseStatus(w, req, err)
		executeTemplate(w, "upload-not-found.tmpl", nil)
		return
	}
	if err != nil {
		errorResponse(w, req, err)
		return
	}

	executeTemplate(w, "upload.tmpl", nil)
}

/* create template */

func handleCreate(w http.ResponseWriter, req *http.Request) {
	executeTemplate(w, "create.tmpl", nil)
}

/* create form */

type createReq struct {
	credential
	Expires time.Duration
}

func handleCreateForm(w http.ResponseWriter, req *http.Request) {
	cred := newCredential(
		req.PostFormValue("Endpoint"),
		req.PostFormValue("Region"),
		req.PostFormValue("AccessKey"),
		req.PostFormValue("SecretKey"),
		req.PostFormValue("Prefix"),
		req.PostFormValue("ACL"),
	)
	if err := cred.validate(); err != nil {
		errorResponse(w, req, err)
		return
	}

	expiresN, err := strconv.ParseUint(req.PostFormValue("Expires"), 10, 64)
	if err != nil {
		errorResponse(w, req, fmt.Errorf("%w: %s", errBadRequest, err))
		return
	}
	expires := time.Duration(expiresN)
	if expires < 10*time.Minute || expires > 90*24*time.Hour {
		errorResponse(w, req, fmt.Errorf("%w: time must be between 10 minutes and 90 days", errBadRequest))
		return
	}

	id := gonanoid.MustGenerate(idAlphabet, 20)

	err = setCredential(id, cred, expires)
	if err != nil {
		errorResponse(w, req, err)
		return
	}

	w.Write([]byte(id))
}
