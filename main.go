package main

import (
	_ "embed"

	"github.com/InditexTech/gh-sherpa/cmd"
)

//go:embed version
var version string

func main() {
	if version == "" {
		version = "Development Build"
	}

	cmd.SetVersion(version)
	cmd.Execute()
}
