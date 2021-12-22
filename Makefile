all: clean test build

jenkins: clean test build

update_dep:
	go get $(DEP)
	go mod tidy

update_all_deps:
	go get -u
	go mod tidy

test:
	go vet ./...
	go test -tags integration ./...

build:
	go build -a -installsuffix cgo examples/consumer/consumer.go
	go build -a -installsuffix cgo examples/publisher/publisher.go

clean:
	rm -rf consumer
	rm -rf publisher


.PHONY: all jenkins update_dep update_all_deps test build clean
