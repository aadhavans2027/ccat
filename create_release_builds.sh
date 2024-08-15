#!/bin/bash

set -euo pipefail

POSSIBLE_GOOS=( "linux" "darwin" )
POSSIBLE_GOARCH=( "amd64" "arm64" )

for OS in "${POSSIBLE_GOOS[@]}"; do
	for ARCH in "${POSSIBLE_GOARCH[@]}"; do
		FOLDER_NAME="ccat-$OS-$ARCH"
		mkdir "${FOLDER_NAME}"
		GOOS=$OS GOARCH=$ARCH go build -o "${FOLDER_NAME}/"
		zip -r "${FOLDER_NAME}" "${FOLDER_NAME}"
		rm -r "${FOLDER_NAME}"
	done
done

