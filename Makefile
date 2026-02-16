GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

# Because we mandate Git Bash, standard Unix find works perfectly
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)

# ==========================================
# Boilerplate Automation Magic
# ==========================================

# 1. Dynamically extract the application name from go.mod
# (e.g., "module github.com/org/loan_service" becomes "loan_service")
APP_NAME := $(shell grep -m 1 '^module ' go.mod | awk '{print $$2}' | awk -F '/' '{print $$NF}')

# 2. Automatically load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# 3. Setup Database variables with fallbacks
# It looks for APP_DATA_DATABASE_SOURCE from your .env. If missing, uses the default.
DB_DSN ?= $(if $(APP_DATA_DATABASE_SOURCE),$(APP_DATA_DATABASE_SOURCE),"postgres://root@localhost:26257/defaultdb?sslmode=disable")
MIGRATION_TABLE ?= "$(APP_NAME)_migration_version"

# 4. Setup Docker variables dynamically
DOCKER_IMAGE ?= $(APP_NAME)
DOCKER_TAG ?= latest

# ==========================================
# Kratos Code Generation
# ==========================================

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./api \
	       --go-http_out=paths=source_relative:./api \
	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make api
	make config
	make generate

# ==========================================
# Database & SQLC
# ==========================================

.PHONY: sqlc
# generate SQLC Go code
sqlc:
	sqlc generate

.PHONY: migrate-up
# run Goose migrations up
migrate-up:
	@echo "🚀 Running $(APP_NAME) migrations against $(DB_DSN)"
	GOOSE_TABLE=$(MIGRATION_TABLE) goose -dir internal/data/migrations postgres $(DB_DSN) up

.PHONY: migrate-down
# run Goose migrations down
migrate-down:
	@echo "⚠️ Reverting $(APP_NAME) migrations"
	GOOSE_TABLE=$(MIGRATION_TABLE) goose -dir internal/data/migrations postgres $(DB_DSN) down

.PHONY: migrate-status
# check Goose migrations status
migrate-status:
	GOOSE_TABLE=$(MIGRATION_TABLE) goose -dir internal/data/migrations postgres $(DB_DSN) status

.PHONY: migrate-create
# create a new Goose migration file. Usage: make migrate-create name=add_users
migrate-create:
	goose -dir internal/data/migrations create $(name) sql

# ==========================================
# Docker Automation
# ==========================================

.PHONY: docker-build
# Build the Debian Docker image
docker-build:
	@echo "🐳 Building Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: docker-run
# Run the container locally (exposing ports 8000 and 9000)
docker-run:
	docker run --rm -p 8000:8000 -p 9000:9000 $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-clean
# Remove dangling/untagged images to free up disk space
docker-clean:
	docker image prune -f

# ==========================================
# Utilities
# ==========================================

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help