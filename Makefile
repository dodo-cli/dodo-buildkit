all: clean test build

.PHONY: clean
clean:
	rm -f dodo-build_*

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run --enable-all

.PHONY: test
test: pkg/types/build_types.pb.go
	go test -cover ./...

.PHONY: build
build: pkg/types/build_types.pb.go
	gox -arch="amd64" -os="darwin linux" ./...

%.pb.go: %.proto
	protoc --go_out=plugins=grpc:. --go_opt=module=github.com/dodo-cli/dodo-build $<
