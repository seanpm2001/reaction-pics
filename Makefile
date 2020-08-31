export ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: test

clean:
	rm reaction-pics
	rm server/static/app.js

node_modules:
	npm install

server/static/app.js: node_modules
	npm run minify

bins: server/static/app.js
	go build

bin/hadolint:
	curl -sL https://github.com/hadolint/hadolint/releases/download/v1.17.3/hadolint-Linux-x86_64 > bin/hadolint && chmod +x bin/hadolint

test: bins bin/hadolint
	./bin/test.sh
	npm test
	git ls-files | grep -e \.sh$ | xargs shellcheck --exclude=SC1091
	bin/hadolint Dockerfile --ignore=DL3008 --ignore=DL4006 --ignore=SC2046 --ignore=SC2006

serve: bins
	./reaction-pics

lsremote:
	bin/rclone --config=bin/rclone.conf ls backblaze:

backupremote:
	bin/rclone --config=bin/rclone.conf sync backblaze: ~/storage/backup-github/reaction-pics-images
