build-all:
	$(MAKE) GOOS=darwin GOARCH=amd64 build
	$(MAKE) GOOS=darwin GOARCH=arm64 build
	$(MAKE) GOOS=linux GOARCH=amd64 build
	$(MAKE) GOOS=linux GOARCH=arm64 build
	$(MAKE) GOOS=windows GOARCH=amd64 build

build:
	CGO_ENABLED=0 go build \
		-ldflags "-s -w" \
		-trimpath \
		-o bin/qbee-cli-$(VERSION).$(GOOS)-$(GOARCH)

test:
	go test ./...