package echoswagger

import (
	"embed"
	"io/fs"
)

var (
	//go:embed dist/*
	dist embed.FS

	// FS holds embedded echo-swagger ui files.
	FS fs.FS
)

func init() {
	var err error
	FS, err = fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
}
