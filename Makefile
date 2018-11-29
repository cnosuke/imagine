NAME     := imagine
VERSION  := $(shell git describe --tags 2>/dev/null)
REVISION := $(shell git rev-parse --short HEAD 2>/dev/null)
SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

export GO111MODULE = on

bin/$(NAME): $(SRCS)
	go build -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME) main.go

.PHONY: test deps

deps:
	go mod vendor

test:
	go test -v ./...
