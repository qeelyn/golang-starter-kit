# Binary names
BUILD_NAME=serve
GO=$(shell which go)

.PHONY: dep
dep:
	dep ensure

.PHONY: test
test:
	$(GO) test -json ./...

.PHONY: build
build: pre-build build-linux

pre-build:
	cp -rf gateway/public cmd/

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o cmd/$(BUILD_NAME) cmd/main.go

build-window:
	GOOS=window GOARCH=amd64 $(GO) build -o cmd/$(BUILD_NAME) cmd/main.go

.PHONY: clean
clean:
	rm -f cmd/$(BUILD_NAME)
	rm -rf cmd/public
	rm -rf cmd/runtime/*.*

.PHONY: docker
docker: clean build
	docker build -t golang-starter-kit .
