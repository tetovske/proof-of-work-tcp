.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: check
check: test lint vet

.PHONY: build
build:
	docker-compose build

.PHONY: up
up:
	docker-compose up

.PHONY: redo
redo: build up
