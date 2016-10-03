all:
	go build
	./test.sh

serve:
	go build
	./devops-reactions-index
