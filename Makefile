all: clean test build

jenkins: install_dep_tool install_go_linter production_restore_deps clean test build

install_dep_tool:
	go get github.com/tools/godep

install_go_linter:
	go get -u -v golang.org/x/lint/golint

initialize_deps:
	go get -d -v ./...
	go get -d -v github.com/stretchr/testify/assert
	go get -v golang.org/x/lint/golint
	godep save

update_deps:
	godep go get -d -v ./...
	godep go get -d -v github.com/stretchr/testify/assert
	godep go get -v golang.org/x/lint/golint
	godep update ./...

test:
	golint ./...
	godep go vet ./...
	godep go test -tags integration ./...

build:
	godep go build -a -installsuffix cgo examples/consumer/consumer.go
	godep go build -a -installsuffix cgo examples/publisher/publisher.go

clean:
	rm -rf consumer
	rm -rf publisher

production_restore_deps:
	godep restore

.PHONY: all jenkins deps install_dep_tool install_go_linter initialize_deps update_deps test build clean production_restore_deps