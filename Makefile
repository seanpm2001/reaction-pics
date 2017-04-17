ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export SERVER_DIR:=$(ROOT_DIR)/server/

all:
	go build
	./test.sh

serve:
	go build
	./reaction-pics
