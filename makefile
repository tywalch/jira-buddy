BINARY_NAME=jira-buddy

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...

## run: run the application
.PHONY: run
run:
	go run main.go

## install: fetch dependencies
.PHONY: install
install:
	go mod download

## clean: remove dead code and binaries
.PHONY: clean
clean: tidy
	go clean
	rm bin/${BINARY_NAME}-darwin
	rm bin/${BINARY_NAME}-linux
	rm bin/${BINARY_NAME}-windows
	rm bin/jira-buddy

## push: push changes to the remote Git repository
.PHONY: push
push: clean audit
	git push

## init: locally initialize the application
.PHONY: init
init: install compile/local

## use: locally initialize and run the application
.PHONY: use
use: init run

## compile: compile binaries for all targets
.PHONY: compile
compile: install compile/darwin compile/linux compile/windows

## compile/local: build the application for use by the current platform
.PHONY: compile/local
compile/local:
	go build -o bin/${BINARY_NAME} main.go

## compile/darwin: compile a binary for mac
.PHONY: compile/darwin
compile/darwin:
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin main.go

## compile/linux: compile a binary for linux
.PHONY: compile/linux
compile/linux:
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux main.go

## compile/windows: compile a binary for windows
.PHONY: compile/windows
compile/windows:
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows main.go

