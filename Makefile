BIN := "./bin/dbshaker"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/dbshaker/main.go

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.49.0

lint: install-lint-deps
	golangci-lint run ./...

#The `-d` flag is not used here in order to make logging more informative
run-integration-test-pg:
	@docker-compose -f tests/docker-compose.yml -f tests/postgres/docker-compose.pg.local.yml \
		--project-directory tests/postgres up \
		--build --abort-on-container-exit

	@docker-compose -f tests/docker-compose.yml -f tests/postgres/docker-compose.pg.local.yml \
		--project-directory tests/postgres down \
		 --rmi local \
		--volumes \
		--timeout 60;