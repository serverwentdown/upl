package main

import (
	"embed"
)

//go:embed web/*.tmpl web/assets/*
var assets embed.FS
