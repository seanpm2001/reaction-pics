#!/bin/bash

cd -P `pwd`
godirs=$(glide novendor)
go test -v $godirs
gofiles=$(find . -name "*.go" | grep -v ./vendor)

VET_ERRORS=$(go vet .)
if [ -n "$VET_ERRORS" ]; then
    echo "Lint failures on:"
    echo "$VET_ERRORS"
    exit 1
fi

fmt_errors=""
for gofile in $gofiles; do
    fmt_errors+=$(gofmt -e -l -d $gofile)
done
if [ -n "$fmt_errors" ]; then
    echo "Fmt failures on:"
    echo "$fmt_errors"
    exit 1
fi

go get -u github.com/golang/lint/golint
LINT_ERRORS=$(golint .)
if [ -n "$LINT_ERRORS" ]; then
    echo "Lint failures on:"
    echo "$LINT_ERRORS"
    exit 1
fi
