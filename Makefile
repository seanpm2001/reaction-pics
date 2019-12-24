export ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: test

vendor:
	dep ensure

node_modules:
	npm install

bins: vendor node_modules
	go build
	npm run minify

bin/hadolint:
	curl -sL https://github.com/hadolint/hadolint/releases/download/v1.17.3/hadolint-Linux-x86_64 > bin/hadolint && chmod +x bin/hadolint

test: bins bin/hadolint
	./bin/test.sh
	npm test
	git ls-files | grep -e \.sh$ | xargs shellcheck
	bin/hadolint Dockerfile --ignore=DL3008 --ignore=DL3009 --ignore=DL4006 --ignore=SC2046 --ignore=SC2006

serve: bins
	./reaction-pics
