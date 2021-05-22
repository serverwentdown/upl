package main

import (
	"encoding/json"
	"errors"
	"html/template"
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

/* templates */

var tmpl = template.Must(template.ParseFS(assets, "web/*.tmpl"))

/* upload template */

func handleUpload(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	_, err := getCredential(vars["id"])
	if errors.Is(err, errNotFound) {
		errorResponseStatus(w, req, err)
		tmpl.ExecuteTemplate(w, "upload-not-found.tmpl", nil)
		return
	}
	if err != nil {
		errorResponse(w, req, err)
		return
	}

	tmpl.ExecuteTemplate(w, "upload.tmpl", nil)
}

/* create template */

func handleCreate(w http.ResponseWriter, req *http.Request) {
	tmpl.ExecuteTemplate(w, "create.tmpl", nil)
}
