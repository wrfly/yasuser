NAME = yasuser

VERSION := $(shell cat VERSION)
COMMITID := $(shell git rev-parse --short HEAD)
BUILDAT := $(shell date +%Y-%m-%d)

CTIMEVAR = -X main.CommitID=$(COMMITID) \
        -X main.Version=$(VERSION) \
        -X main.BuildAt=$(BUILDAT)
GO_LDFLAGS = -ldflags "-w $(CTIMEVAR) -s"

.PHONY: build
build:
	go build -tags "$(BUILDTAGS)" $(GO_LDFLAGS) -o $(NAME) .

.PHONY: test
test:
	go test -v --cover `glide nv`

.PHONY: dev
dev: asset build
	rm -f $(NAME).db
	YASUSER_DEBUG=true YASUSER_SERVER_PPROF=true ./$(NAME)

.PHONY: pprof
pprof: asset build
	rm -f $(NAME).db
	YASUSER_SERVER_PPROF=true ./$(NAME)

.PHONY: img
img:
	docker build -t wrfly/$(NAME):$(VERSION) -t wrfly/$(NAME) -t wrfly/$(NAME):develop .

.PHONY: push-img
push-img:
	docker push wrfly/$(NAME)
	docker push wrfly/$(NAME):$(VERSION)

.PHONY: push-dev-img
push-dev-img:
	docker push wrfly/$(NAME):develop

.PHONY: tools
tools:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

.PHONY: asset
asset:
	go-bindata-assetfs -nometadata -prefix routes/index -pkg routes routes/index/...
	mv bindata_assetfs.go routes/asset.go
	gofmt -w routes/asset.go
