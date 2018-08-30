#!/bin/bash

set -e

cd -P `pwd`
echo "" > coverage.txt
for godir in $(go list ./... | grep -v vendor); do
    go test -coverprofile=coverage.out $godir -covermode=atomic
    if [ -f coverage.out ]
    then
        cat coverage.out | grep -v "mode: set" >> coverage.txt
    fi
done
rm coverage.out


govet_errors=$(go vet ./...)
if [ -n "$govet_errors" ]; then
    echo "Vet failures on:"
    echo "$govet_errors"
    exit 1
fi

gofiles=$(find . -name "*.go" | grep -v ./vendor)

fmt_errors=$(gofmt -e -l -d $gofiles)
if [ -n "$fmt_errors" ]; then
    echo "Fmt failures on:"
    echo "$fmt_errors"
    exit 1
fi

go get -u github.com/golang/lint/golint
golint_errors=""
for gofile in $gofiles; do
    golint_errors+=$(golint $gofile)
done
if [ -n "$golint_errors" ]; then
    echo "Lint failures on:"
    echo "$golint_errors"
    exit 1
fi
