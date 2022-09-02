all: test

clean:
	rm reaction-pics
	rm server/static/app.js

bins:
	go build -race .

test: bins
	./bin/test.sh
	npm test
	git ls-files | grep -e \.sh$ | xargs shellcheck --exclude=SC1091

serve: bins
	./reaction-pics

lsremote:
	bin/rclone --config=bin/rclone.conf ls backblaze:

backupremote:
	bin/rclone --config=bin/rclone.conf sync backblaze: ~/storage/backup-github/reaction-pics-images
