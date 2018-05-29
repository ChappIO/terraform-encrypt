#!/bin/bash
set -e

CGO_ENABLED=0
for GOARCH in amd64 386
do
    for GOOS in darwin windows linux
    do
        OUTPUT="build/terraform-encrypt-${GOOS}_${GOARCH}"
        if [ $GOOS = "windows" ]
        then
            OUTPUT="${OUTPUT}.exe"
        fi

        export GOARCH
        export CGO_ENABLED
        export GOOS
        echo "Building ${OUTPUT}..."
        go build -o ${OUTPUT} main.go
    done
done