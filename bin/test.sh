#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd -P "$DIR/.." || exit 1
pwd

go test -race -coverprofile=coverage.txt -covermode=atomic ./...
go vet ./...

gofmt_errors=$(gofmt -e -l -d -s .)
if [ -n "$gofmt_errors" ]; then
    echo "Fmt failures on:"
    echo "$gofmt_errors"
    exit 1
fi

go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck -checks all ./...

go mod tidy
gosumdiff="$(git diff go.sum)"
if [ -n "$gosumdiff" ]; then
    echo "Go.sum not up to date:"
    echo "$gosumdiff"
    exit 1
fi
