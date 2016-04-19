all: update_deps build test

deps:
	go get -d -v ./...
	go get -d -v github.com/stretchr/testify/assert
	go get -v github.com/golang/lint/golint

update_deps:
	go get -d -v -u ./...
	go get -d -v -u github.com/stretchr/testify/assert
	go get -v -u github.com/golang/lint/golint


test:
	golint ./...
	go vet ./...
	go test -tags integration ./...

build:
	go build examples/consumer/consumer.go
	go build examples/publisher/publisher.go

.PHONY: deps update_deps test build
