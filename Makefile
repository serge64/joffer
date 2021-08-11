.PHONE: build
build:
	go build -v ./cmd/app/

.PHONE: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build