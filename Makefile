VERSION=$(shell git describe --tags --always --dirty)

build-local:
	CGO_ENABLED=0 go build \
		-ldflags "-s -w -X go.qbee.io/client.Version=$(VERSION)" \
		-trimpath \
		-o bin/qbee-cli \
		./cmd

build-all:
	$(MAKE) GOOS=darwin GOARCH=amd64 build
	$(MAKE) GOOS=darwin GOARCH=arm64 build
	$(MAKE) GOOS=linux GOARCH=amd64 build
	$(MAKE) GOOS=linux GOARCH=arm64 build
	$(MAKE) GOOS=windows GOARCH=amd64 build

build:
	CGO_ENABLED=0 go build \
		-ldflags "-s -w -X go.qbee.io/client.Version=$(VERSION)" \
		-trimpath \
		-o bin/qbee-cli-$(VERSION).$(GOOS)-$(GOARCH) \
		./cmd

test:
	go test ./...

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.64.5 golangci-lint run
