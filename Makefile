# Common
NAME ?= joffer
VERSION ?= $(shell git tag --sort -version:refname | head -1)
POSTGRES_URL ?= postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

# Docker-compose
DOCKER_COMPOSE_FILE ?= deployments/docker-compose/docker-compose.yml

# Build
BUILD_CMD ?= CGO_ENABLED=0 go build -o bin/${NAME} -ldflags '-v -w -s' ./cmd/${NAME}
DOCKER_COMPOSE_UP ?= docker-compose -f ${DOCKER_COMPOSE_FILE} up -d
DOCKER_COMPOSE_DOWN ?= docker-compose -f ${DOCKER_COMPOSE_FILE} down
MIGRATIONS_UP ?= migrate -database ${POSTGRES_URL} -path migrations/ up
MIGRATIONS_DOWN ?= migrate -database ${POSTGRES_URL} -path migrations/ down

.PHONY: docker_build
docker_build:
	docker build -t joffer:${VERSION} ${DOCKER_APP_FILENAME}

.PHONY: compose_up
compose_up:
	${DOCKER_COMPOSE_UP}

.PHONY: compose_down 
compose_down:
	${DOCKER_COMPOSE_DOWN}

.PHONY: migrate_up
migrate_up:
	${MIGRATIONS_UP}

.PHONY: migrate_down
migrate_down:
	${MIGRATIONS_DOWN}

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: build
build:
	${BUILD_CMD}
	cp .env bin/

.DEFAULT_GOAL := build
