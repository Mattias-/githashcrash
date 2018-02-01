#!/usr/bin/env bash
set -eu

docker run --rm \
    -v "$PWD":/go/src/githashcrash \
    -w /go/src/githashcrash \
    -e "GOOS=darwin" \
    golang:1.9.3-alpine3.7 \
    go build -o bin/githashcrash-macos

docker run --rm \
    -v "$PWD":/go/src/githashcrash \
    -w /go/src/githashcrash \
    -e "GOOS=linux" \
    golang:1.9.3-alpine3.7 \
    go build -o bin/githashcrash

docker build . -t githashcrash
