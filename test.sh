#!/bin/bash

cd -P `pwd`
godirs=$(glide novendor)
go test -v $godirs
gofiles=$(find . -name "*.go" | grep -v ./vendor)

govet_errors=""
for gofile in $gofiles; do
    govet_errors+=$(go vet $gofile 2>&1)
done
if [ -n "$govet_errors" ]; then
    echo "Vet failures on:"
    echo "$govet_errors"
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
golint_errors=""
for gofile in $gofiles; do
    golint_errors+=$(golint $gofile)
done
if [ -n "$golint_errors" ]; then
    echo "Lint failures on:"
    echo "$golint_errors"
    exit 1
fi
