BIN="./bin"
NAME="envc"
SRC=$(shell find . -name "*.go")

.PHONY: fmt lint test install_deps clean

default: all

all: fmt test

fmt:
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

test: install_deps
	$(info ******************** running tests ********************)
	go test -v ./...

install_deps:
	$(info ******************** downloading dependencies ********************)
	go get -v ./...

build:
	$(info ******************** building ********************)
	go build -o $(BIN)/$(NAME)

clean:
	rm -rf $(BIN)