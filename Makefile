export ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: test

bins:
	go build

test: bins
	./bin/test.sh

serve: bins
	./reaction-pics

clear-cache:
	rm -f tumblr/data/static/*
