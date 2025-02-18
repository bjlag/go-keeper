BUILD_VERSION = "v1.0.0"
BUILD_DATE = $(shell date +'%Y/%m/%d %H:%M:%S')
BUILD_COMMIT = $(shell git rev-parse --short HEAD)

NAME=$(shell basename "$(PWD)")
DIR=$(shell pwd)

## lint: start linter
lint:
	@echo " > Start linter"
	@golangci-lint run

## fmt: start fmt
fmt:
	@echo " > Start fmt"
	@goimports -local "github.com/bjlag/go-keeper" -d -w $$(find . -type f -name '*.go' -not -path "*_mock.go")

## test: start testing
test:
	@echo " > Testing"
	@go test -v $(DIR)/...

## tidy: start `go mod tidy`
tidy:
	@echo "  >  Go mod tidy"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod tidy

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo