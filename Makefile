.PHONY: build clean tool lint help

all: build

build:
	@echo "App is creating. Please wait ..."
	@go build -o novels-spider -v . 
	@echo "App is created"

tool:
	go vet ./...; true
	gofmt -w .

lint:
	golint ./...

clean:
	rm -rf novels-spider
	go clean -i .

test:
	@echo "Test --- START"
	@go test -v pkg/file/*.go
	@go test -v pkg/queue/*.go
	@echo "Test --- END"


help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
