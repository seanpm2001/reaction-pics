export ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: test

test:
	go build
	./bin/test.sh

serve:
	go build
	./reaction-pics
