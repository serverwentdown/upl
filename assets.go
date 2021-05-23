package main

import (
	"embed"
	"io/fs"
)

//go:embed web/*.tmpl web/assets/*
var assets embed.FS

var assetsWeb = fsMust(fs.Sub(assets, "web"))

func fsMust(fs fs.FS, err error) fs.FS {
	if err != nil {
		panic(err)
	}
	return fs
}
