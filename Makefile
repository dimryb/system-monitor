BIN := "./bin/system-monitor"

GIT_HASH := $(shell git log --format="%h" -n 1)

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    DATE_CMD = date -u +'%Y-%m-%dT%H:%M:%S'
    GO_PATH := $(shell go env GOPATH)
else #windows
    DATE_CMD = powershell.exe -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ss'"
    GO_PATH := $(shell go env GOPATH | tr '\\' '/')
endif

LDFLAGS := -X main.release="develop" \
    -X main.buildDate=$(shell $(DATE_CMD)) \
    -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/monitor

run: build
	$(BIN) -config ./configs/monitor.yaml

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
	sh -s -- -b $(GO_PATH)/bin v1.64.8

lint: install-lint-deps
	golangci-lint run --config golangci.yml ./...

generate:
	protoc \
		-I proto \
		--go_out=proto --go_opt=paths=source_relative \
		--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
		proto/monitor/*.proto

build-img:
	docker build \
        --build-arg LDFLAGS="$(LDFLAGS)" \
        -t $(DOCKER_IMG) \
        -f $(DOCKERFILE_PATH) .

build-linux-img:
	$(MAKE) build-img \
            DOCKER_IMG=system-monitor:develop \
            DOCKERFILE_PATH=build/linux/Dockerfile

run-linux-img:
	docker run --rm system-monitor:develop