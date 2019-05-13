#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

cd -P "$(pwd)"
echo "" > coverage.txt
for godir in $(go list ./... | grep -v vendor); do
    go test -coverprofile=coverage.out "$godir" -covermode=atomic
    if [ -f coverage.out ]
    then
        grep -v "mode: set" coverage.out >> coverage.txt
    fi
done
rm coverage.out

govet_errors=$(go vet ./...)
if [ -n "$govet_errors" ]; then
    echo "Vet failures on:"
    echo "$govet_errors"
    exit 1
fi

gofiles=$(git ls-files | grep -F .go)

gofmt_errors=""
for gofile in $gofiles; do
    gofmt_errors+=$(gofmt -e -l -d "$gofile")
done
if [ -n "$gofmt_errors" ]; then
    echo "Fmt failures on:"
    echo "$gofmt_errors"
    exit 1
fi

go get -u golang.org/x/lint/golint
golint_errors=""
for gofile in $gofiles; do
    golint_errors+=$(golint "$gofile")
done
if [ -n "$golint_errors" ]; then
    echo "Lint failures on:"
    echo "$golint_errors"
    exit 1
fi
