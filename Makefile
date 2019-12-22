export ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: test

bins:
	go build
	npm install

test: bins
	./bin/test.sh
	npm test
	git ls-files | grep -e \.sh$ | xargs shellcheck

serve: bins
	npm run minify
	./reaction-pics
