ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export SERVER_TEMPLATES:=$(ROOT_DIR)/server/templates/

all:
	go build
	./test.sh

serve:
	go build
	./devops-reactions-index
