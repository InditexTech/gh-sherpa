// SPDX-FileCopyrightText: 2023 INDITEX S.A
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	_ "embed"
	"strings"

	"github.com/InditexTech/gh-sherpa/cmd"
	"github.com/InditexTech/gh-sherpa/pkg/metadata"
)

//go:embed version
var version string

func main() {
	if version == "" {
		version = "Development Build"
	}

	metadata.Version = strings.TrimSpace(version)

	cmd.SetVersion(version)
	cmd.Execute()
}
