#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd -P "$DIR/.." || exit 1
pwd

go test -coverprofile=coverage.txt -covermode=atomic ./...
go vet ./...

gofiles=$(git ls-files | grep -F .go)

gofmt_errors=""
for gofile in $gofiles; do
    gofmt_errors+=$(gofmt -e -l -d -s "$gofile")
done
if [ -n "$gofmt_errors" ]; then
    echo "Fmt failures on:"
    echo "$gofmt_errors"
    exit 1
fi

go install golang.org/x/lint/golint@latest
golint_errors=""
for gofile in $gofiles; do
    golint_errors+=$(golint "$gofile")
done
if [ -n "$golint_errors" ]; then
    echo "Lint failures on:"
    echo "$golint_errors"
    exit 1
fi

go mod tidy
gosumdiff="$(git diff go.sum)"
if [ -n "$gosumdiff" ]; then
    echo "Go.sum not up to date:"
    echo "$gosumdiff"
    exit 1
fi
