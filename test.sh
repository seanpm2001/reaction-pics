#!/bin/bash

set -e

cd -P `pwd`
go test -v $(./glide novendor)

rm -rf vendor

VET_ERRORS=$(go vet .)
if [ -n "$VET_ERRORS" ]; then
    echo "Lint failures on:"
    echo "$VET_ERRORS"
    ./glide install
    exit 1
fi

FMT_ERRORS=$(gofmt -e -l -d $(./glide novendor))
if [ -n "$FMT_ERRORS" ]; then
    echo "Fmt failures on:"
    echo "$FMT_ERRORS"
    ./glide install
    exit 1
fi

go get -u github.com/golang/lint/golint
LINT_ERRORS=$(golint $(./glide novendor))
if [ -n "$LINT_ERRORS" ]; then
    echo "Lint failures on:"
    echo "$LINT_ERRORS"
    ./glide install
    exit 1
fi

./glide install
