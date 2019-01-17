NAME     := imagine
VERSION  := $(shell git describe --tags 2>/dev/null)
REVISION := $(shell git rev-parse --short HEAD 2>/dev/null)
SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

export GO111MODULE = on

bin/$(NAME): $(SRCS)
	go build -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME) main.go

.PHONY: test deps build-for-docker build-docker build-static copy-static-from-client

build-for-docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$(NAME)

deps:
	go mod vendor

test:
	go test -v ./...

build-static:
	cd ./_front && yarn build

copy-static-from-client:
	rm -r ./static && cp -r ./_front/dist/ ./static

build-docker: build-for-docker build-static copy-static-from-client
	docker build -t cnosuke/imagine:latest
