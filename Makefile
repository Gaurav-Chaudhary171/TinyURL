.PHONY: all build run clean proto

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

# Binary name
BINARY_NAME=tinyurl

# Protobuf parameters
PROTOC=protoc
PROTO_DIR=.
GO_OUT_DIR=proto

all: proto build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

run: build
	./$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(GO_OUT_DIR)/*.pb.go

proto:
	$(PROTOC) -I. -I./googleapis \
		--go_out=./proto --go_opt=paths=source_relative \
		--go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
		proto/tinyurl.proto

deps:
	$(GOGET) -u google.golang.org/protobuf/cmd/protoc-gen-go
	$(GOGET) -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

test:
	$(GOTEST) -v ./...

help:
	@echo "Available targets:"
	@echo "  all        - Generate protobuf and build the project"
	@echo "  build      - Build the project"
	@echo "  run        - Build and run the project"
	@echo "  clean      - Clean build artifacts"
	@echo "  proto      - Generate protobuf code"
	@echo "  deps       - Install required dependencies"
	@echo "  test       - Run tests" 