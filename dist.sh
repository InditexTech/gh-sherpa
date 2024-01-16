#!/bin/bash
set -e

NAME="gh-sherpa"
ROOT_DIR="$(dirname "$0")"
VERSION=$(head -1 "${ROOT_DIR}/version")
DIST_DIR="${ROOT_DIR}/dist"

if [ "${VERSION}" == "" ]; then
  echo "Version is required."
  exit 1
fi

# Build distributable files of the extension.
rm -rf "${DIST_DIR}"
echo "==> Building distributables for the supported architectures..."

# Disable CGO when building.
export CGO_ENABLED=0

# MacOS
echo "-> MacOS x86_64"
GOOS=darwin GOARCH=amd64  go build -o "dist/${NAME}-darwin-x86_64-${VERSION}"
echo "-> MacOS arm64"
GOOS=darwin GOARCH=arm64  go build -o "dist/${NAME}-darwin-arm64-${VERSION}"

#Linux
echo "-> Linux i386"
GOOS=linux GOARCH=386     go build -o "dist/${NAME}-linux-i386-${VERSION}"
echo "-> Linux x86_64"
GOOS=linux GOARCH=amd64   go build -o "dist/${NAME}-linux-x86_64-${VERSION}"
echo "-> Linux arm64"
GOOS=linux GOARCH=arm64   go build -o "dist/${NAME}-linux-arm64-${VERSION}"

# Windows
echo "-> Windows i386"
GOOS=windows GOARCH=386   go build -o "dist/${NAME}-windows-i386-${VERSION}.exe"
echo "-> Windows x86_64"
GOOS=windows GOARCH=amd64 go build -o "dist/${NAME}-windows-x86_64-${VERSION}.exe"
echo "-> Windows arm64"
GOOS=windows GOARCH=amd64 go build -o "dist/${NAME}-windows-arm64-${VERSION}.exe"

echo "==> Done!"
