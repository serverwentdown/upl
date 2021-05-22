package main

import (
	"embed"
)

//go:embed web/index.html web/assets/*
var assets embed.FS
