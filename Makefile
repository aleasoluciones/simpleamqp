all: clean test build

update_dep:
	go get $(DEP)
	go mod tidy

update_all_deps:
	go get -u all
	go mod tidy

test:
	go vet ./...
	go clean -testcache
	go test -v -tags integration ./... -timeout 60s

build:
	go build examples/consumer/consumer.go
	go build examples/publisher/publisher.go

clean:
	rm -f consumer
	rm -f publisher

start_dependencies:
	docker-compose -f dev/simpleamqp_devdocker/docker-compose.yml up -d

stop_dependencies:
	docker-compose -f dev/simpleamqp_devdocker/docker-compose.yml stop

rm_dependencies:
	docker-compose -f dev/simpleamqp_devdocker/docker-compose.yml down -v


.PHONY: all update_dep update_all_deps test build clean start_dependencies stop_dependencies rm_dependencies
