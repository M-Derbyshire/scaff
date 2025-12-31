.DEFAULT_GOAL := build-dev

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

test-unit:
	./scripts/test_unit.sh
.PHONY:test-unit

test-e2e:
	./scripts/test_e2e.sh
.PHONY:test-e2e

build-dev: vet lint
	go build .
.PHONY:build-dev

build-prod: vet lint
	go build -ldflags=-w .
.PHONY:build-prod