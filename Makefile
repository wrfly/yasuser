.PHONY: build test

NAME = short-url

VERSION := $(shell cat VERSION)
COMMITID := $(shell git rev-parse --short HEAD)
BUILDAT := $(shell date +%Y-%m-%d)

CTIMEVAR = -X main.CommitID=$(COMMITID) \
        -X main.Version=$(VERSION) \
        -X main.BuildAt=$(BUILDAT)
GO_LDFLAGS = -ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC = -ldflags "-w $(CTIMEVAR) -extldflags -static"

build:
	go build -tags "$(BUILDTAGS)" $(GO_LDFLAGS) -o $(NAME) .

test:
	go test --cover .

dev: build
	./$(NAME) -d

img:
	docker build -t wrfly/short-url .

push-img:
	docker push wrfly/short-url