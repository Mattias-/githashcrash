#!/usr/bin/env bash
set -eu

docker run --rm \
    -v "$PWD":/go/src/githashcrash \
    -w /go/src/githashcrash \
    -e "GOOS=darwin" \
    golang:1.8 \
    go build -o bin/githashcrash-macos

docker run --rm \
    -v "$PWD":/go/src/githashcrash \
    -w /go/src/githashcrash \
    -e "GOOS=linux" \
    golang:1.8 \
    go build -o bin/githashcrash

docker build . -t githashcrash
