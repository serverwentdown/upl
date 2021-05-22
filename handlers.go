package main

import (
	"html/template"
	"net/http"
	"os"
)

var globalStore store

func setupHandlers() {
	var err error
	globalStore, err = newRedisStore(os.Getenv("REDIS_CONNECTION"))
	if err != nil {
		panic(err)
	}
}

/* templates */

var tmpl = template.Must(template.ParseFS(assets, "web/*.tmpl"))

/* upload template */

func handleUpload(w http.ResponseWriter, req *http.Request) {
	tmpl.ExecuteTemplate(w, "upload.tmpl", nil)
}
