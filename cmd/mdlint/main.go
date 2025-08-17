// Copyright 2025 Sam Caldwell
//
// MdLint is a command-line application for linting Markdown files.
package main

import (
	"flag"
	"fmt"

	"github.com/sam-caldwell/mdlint/internal/version"
)

// main is the application entry point.
func main() {
	showVersion := flag.Bool("version", false, "Print application version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Println(version.Version)
		return
	}
	// TODO: implement CLI
}
