.PHONY: all test build clean

all: test build

build: 
	mkdir -p build
	go build -o build ./...

test:
	go test -v -coverprofile=tests/results/cover.out ./...

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	go clean ./...
