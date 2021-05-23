package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
