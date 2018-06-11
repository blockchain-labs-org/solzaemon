.DEFAULT_GOAL = test

test:
	go test -v ./...

solzaemon: *.go */*.go
	go build .

run: solzaemon
	./solzaemon
