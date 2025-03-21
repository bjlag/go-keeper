BUILD_VERSION ?= "v1.0.0"
BUILD_DATE ?= $(shell date +'%Y/%m/%d %H:%M:%S')

NAME = $(shell basename "$(PWD)")
DIR = $(shell pwd)
DOCKER_FILE = "docker-compose.local.yaml"

.DEFAULT_GOAL := up

.PHONY: all
all: help

## up: start app
.PHONY: up
up: docker-up wait-db migrate

## down: stop app
.PHONY: down
down: docker-down

wait-db:
	@echo " > Wait DB"
	@sleep 5

## clean: stop app and remove volumes
.PHONY: clean
clean: docker-down-clear

## docker-up: start docker
.PHONY: docker-up
docker-up:
	@echo " > Start docker"
	@docker-compose -f $(DIR)/docker/$(DOCKER_FILE) up -d

## docker-down: stop docker
.PHONY: docker-down
docker-down:
	@echo " > Stop docker"
	@docker-compose -f $(DIR)/docker/$(DOCKER_FILE) down --remove-orphans

## docker-down-clear: stop docker and remove volumes
.PHONY: docker-down-clear
docker-down-clear:
	@echo " > Stop docker and remove volumes"
	@docker-compose -f $(DIR)/docker/$(DOCKER_FILE) down -v --remove-orphans

## migrate: apply migrations
.PHONY: migrate
migrate:
	@echo " > Apply migrations"
	@go run $(DIR)/cmd/migrator -c="./config/migrator_server.yaml"

## lint: start linter
.PHONY: lint
lint:
	@echo " > Start linter"
	@golangci-lint run

## fmt: start fmt
.PHONY: fmt
fmt:
	@echo " > Start fmt"
	@goimports -local "github.com/bjlag/go-keeper" -d -w $$(find . -type f -name '*.go' -not -path "*_mock.go" -not -path "./internal/generated/*")

## test: start testing
.PHONY: test
test:
	@echo " > Testing"
	@go test -v $(DIR)/...

## tidy: start `go mod tidy`
.PHONY: tidy
tidy:
	@echo "  >  Go mod tidy"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod tidy

## proto: generate grpc client/server from proto files
.PHONY: proto
proto:
	@echo "  >  Generate gRPC"
	@protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import --go-grpc_opt=require_unimplemented_servers=false proto/*

## doc: open documentation
.PHONY: doc
doc:
	godoc -http=:8888 -play

build-client:
	@echo " > Building client for OSX"
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)'" -o ./cmd/client/client_darwin_amd64 ./cmd/client/.
	@echo " > Building client for Linux"
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)'" -o ./cmd/client/client_linux_amd64 ./cmd/client/.
	@echo " > Building client for Windows"
	@GOOS=windows GOARCH=amd64 go build -ldflags "-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)'" -o ./cmd/client/client_windows_amd64 ./cmd/client/.

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo