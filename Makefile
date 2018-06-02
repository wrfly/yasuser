.PHONY: build test dev img push-img

NAME = yasuser

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
	YASUSER_DEBUG=true ./$(NAME)

img:
	docker build -t wrfly/$(NAME):$(VERSION) -t wrfly/$(NAME) .

push-img:
	docker push wrfly/$(NAME)
	docker push wrfly/$(NAME):$(VERSION)