package docs

import "embed"

//go:embed api/*
var StaticFiles embed.FS
