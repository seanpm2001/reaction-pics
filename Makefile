export ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export SERVER_DIR:=$(ROOT_DIR)/server/

all: test

test:
	go build
	./bin/test.sh

serve:
	go build
	./reaction-pics
