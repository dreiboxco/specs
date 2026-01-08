package main

import (
	"os"

	"github.com/dreibox/specs/internal/cli"
	"github.com/dreibox/specs/internal/adapters"
)

var (
	version = "dev" // Injetado durante build via ldflags
)

func main() {
	fs := adapters.NewFileSystem()
	router := cli.NewRouter(fs, version)
	
	code := router.Run(os.Args[1:])
	os.Exit(code)
}

