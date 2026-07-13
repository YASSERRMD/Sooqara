.PHONY: build run test fmt vet lint clean docker

APP_NAME := sooqara
BIN_DIR  := bin
MODULE   := github.com/yasserrmd/sooqara

build:
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

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
