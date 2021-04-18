.PHONY: build clean gomodgen format

build: gomodget format
	go build ./githubtest/githubtest.go

# NOTE: Add the `-test.v` flag for verbose logging.
test: build format
	go test -test.v ./githubtest/...
	# go test ./...

clean:
	rm -rf ./bin ./vendor Gopkg.lock

gomodget:
	go get -v all

format:
	go fmt ./...