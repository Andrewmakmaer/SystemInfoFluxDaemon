BIN := "./.bin/daemon"
DAEMON_IMG := "daemon:develop"
CLIENT_IMG := "client:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/

run: build
	$(BIN) -config ./configs/

test:
	go test -race ./internal/...

generate:
	rm -rf ./internal/server/pb
	mkdir -p ./internal/server/pb

	protoc \
		--proto_path=api/ \
		--go_out=internal/server/pb \
		--go-grpc_out=internal/server/pb \
		api/*.proto

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.1

docker:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DAEMON_IMG) \
		-f DockerBuild/DaemonDockerfile .
	
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(CLIENT_IMG) \
		-f DockerBuild/ClientDockerfile .

docker-daemon-run: docker
	docker run -d -e "DAEMON_MODES=cpu la" -e "DAEMON_LOG_LEVEL=INFO" -e "DAEMON_PORT=8765" -p 8765:8765 $(DAEMON_IMG)

integrations: docker-daemon-run
	INTEGRATION_DAEMON_URL=":8765" go run integration_test/main.go

lint: install-lint-deps
	golangci-lint run --fix ./...
	