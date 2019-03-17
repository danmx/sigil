
# Change this and commit to create new release
VERSION=0.0.1

SRC = $(wildcard pkg/*) $(wildcard cmd/*)
#MODULE = github.com/danmx/sigil
NAME = sigil
REVISION := $(shell git rev-parse --short HEAD;)

export GO111MODULE = on

.PHONY: bootstrap
bootstrap:
	@go mod download && go mod vendor

.PHONY: build
build: bootstrap build-linux build-mac build-windows build-docker

.PHONY: build-dev
build: bootstrap build-linux-dev build-mac-dev build-windows-dev build-docker-dev

build-windows: export GOARCH=amd64
build-windows:
	@GOOS=windows go build -mod=vendor -v \
		--ldflags="-w -X main.AppName=$(NAME) -X main.Version=$(VERSION) \
		-X main.Revision=$(REVISION)" -o bin/windows/amd64/$(NAME) cmd/$(NAME)/main.go

build-linux: export GOARCH=amd64
build-linux: export CGO_ENABLED=0
build-linux:
	@GOOS=linux go build -mod=vendor -v \
		--ldflags="-w -X main.AppName=$(NAME) -X main.Version=$(VERSION) \
		-X main.Revision=$(REVISION)" -o bin/linux/amd64/$(NAME) cmd/$(NAME)/main.go

build-mac: export GOARCH=amd64
build-mac: export CGO_ENABLED=0
build-mac:
	@GOOS=darwin go build -mod=vendor -v \
		--ldflags="-w -X main.AppName=$(NAME) -X main.Version=$(VERSION) \
		-X main.Revision=$(REVISION)" -o bin/darwin/amd64/$(NAME) cmd/$(NAME)/main.go

build-docker:
	@docker build -t $(NAME):$(VERSION) .

build-windows-dev: export GOARCH=amd64
build-windows-dev:
	@GOOS=windows go build -mod=vendor -v \
		--ldflags="-w -X main.LogLevel=debug -X main.AppName=$(NAME) \
		-X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/windows/amd64/dev/$(NAME) cmd/$(NAME)/main.go

build-linux-dev: export GOARCH=amd64
build-linux-dev: export CGO_ENABLED=0
build-linux-dev:
	@GOOS=linux go build -mod=vendor -v \
		--ldflags="-w -X main.LogLevel=debug -X main.AppName=$(NAME) \
		-X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/linux/amd64/dev/$(NAME) cmd/$(NAME)/main.go

build-mac-dev: export GOARCH=amd64
build-mac-dev: export CGO_ENABLED=0
build-mac-dev:
	@GOOS=darwin go build -mod=vendor -v \
		--ldflags="-w -X main.LogLevel=debug -X main.AppName=$(NAME) \
		-X main.Version=$(VERSION) -X main.Revision=$(REVISION)" \
		-o bin/darwin/amd64/dev/$(NAME) cmd/$(NAME)/main.go

build-docker-dev:
	@docker build -t $(NAME):$(VERSION)-dev -f Dockerfile-dev .

.PHONY: get-version
get-version:
	@echo $(VERSION)

.PHONY: clean
clean:
	@git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf

.PHONY: test
test:
	@go test -v -coverpkg=$(SRC) -failfast -timeout=2m

.PHONY: tidy
	@go mod tidy
