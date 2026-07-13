.PHONY: build run test fmt vet lint clean docker release release-version

APP_NAME := sooqara
BIN_DIR  := bin
MODULE   := github.com/yasserrmd/sooqara

build:
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

release:
	go build -ldflags "$(shell go run ./internal/release/build_flags.go dev $(COMMIT))" -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

release-version:
	@if [ -z "$(VERSION)" ]; then echo "Usage: make release-version VERSION=v1.0.0 COMMIT=abc123"; exit 1; fi
	go build -ldflags "-X $(MODULE)/internal/version.Version=$(VERSION) -X $(MODULE)/internal/version.Commit=$(COMMIT) -X $(MODULE)/internal/version.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

run:
ifeq ($(wildcard .env),)
	$(error .env file not found. Copy .env.example to .env and set your AGNES_API_KEY.)
endif
	go run ./cmd/$(APP_NAME)

test:
	go test -race -count=1 ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: vet

clean:
	rm -rf $(BIN_DIR)
	find . -name '*.db' -delete

docker:
	docker build -t $(APP_NAME) .
