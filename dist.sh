#!/bin/bash
set -e

NAME="gh-sherpa"
ROOT_DIR="$(dirname "$0")"
DIST_DIR="${ROOT_DIR}/dist"

# Build distributable files of the extension.
rm -rf "${DIST_DIR}"
echo "==> Building distributables for the supported architectures..."

# Disable CGO when building.
export CGO_ENABLED=0

# MacOS
echo "-> MacOS amd64"
GOOS=darwin GOARCH=amd64  go build -o "dist/${NAME}-darwin-amd64"
echo "-> MacOS arm64"
GOOS=darwin GOARCH=arm64  go build -o "dist/${NAME}-darwin-arm64"

#Linux
echo "-> Linux i386"
GOOS=linux GOARCH=386     go build -o "dist/${NAME}-linux-386"
echo "-> Linux amd64"
GOOS=linux GOARCH=amd64   go build -o "dist/${NAME}-linux-amd64"
echo "-> Linux arm64"
GOOS=linux GOARCH=arm64   go build -o "dist/${NAME}-linux-arm64"

# Windows
echo "-> Windows i386"
GOOS=windows GOARCH=386   go build -o "dist/${NAME}-windows-386.exe"
echo "-> Windows amd64"
GOOS=windows GOARCH=amd64 go build -o "dist/${NAME}-windows-amd64.exe"
echo "-> Windows arm64"
GOOS=windows GOARCH=arm64 go build -o "dist/${NAME}-windows-arm64.exe"

echo "==> Done!"
