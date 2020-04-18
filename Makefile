
# Change this and commit to create new release
override VERSION ?= 0.4.1
NAME = sigil
REPO = danmx/$(NAME)
MODULE = github.com/$(REPO)
override REVISION ?= $(shell git rev-parse HEAD;)

export GO111MODULE = on
export CGO_ENABLED=0
export GOARCH=amd64

.PHONY: all
all: bootstrap format lint test buildallnull release

.PHONY: bootstrap
bootstrap:
	go mod download && go mod verify
bootstrap: generate

.PHONY: build
build: bootstrap build-linux build-darwin build-windows

.PHONY: build-dev
build-dev: bootstrap build-linux-dev build-darwin-dev build-windows-dev

.PHONY: release
release: build release-windows release-linux release-darwin

.PHONY: buildallnull
buildallnull:
	GOOS=windows go build -mod=readonly -v \
		--ldflags="-w -s \
			-X $(MODULE)/cmd.LogLevel=panic \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o /dev/null main.go
	GOOS=linux go build -mod=readonly -v \
		--ldflags="-w -s \
			-X $(MODULE)/cmd.LogLevel=panic \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o /dev/null main.go
	GOOS=darwin go build -mod=readonly -v \
		--ldflags="-w -s \
			-X $(MODULE)/cmd.LogLevel=panic \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o /dev/null main.go
buildallnull: build-docker

.PHONY: release-windows
release-windows:
	mkdir -p dist && tar -czvf dist/$(NAME)_windows-$(GOARCH).tar.gz -C bin/release/windows/$(GOARCH)/ $(NAME).exe

.PHONY: release-linux
release-linux:
	mkdir -p dist && tar -czvf dist/$(NAME)_linux-$(GOARCH).tar.gz -C bin/release/linux/$(GOARCH)/ $(NAME)

.PHONY: release-darwin
release-darwin:
	mkdir -p dist && tar -czvf dist/$(NAME)_darwin-$(GOARCH).tar.gz -C bin/release/darwin/$(GOARCH)/ $(NAME)

.PHONY: build-windows
build-windows:
	GOOS=windows go build -mod=readonly -v \
		--ldflags="-w -s \
			-X $(MODULE)/cmd.LogLevel=panic \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/release/windows/$(GOARCH)/$(NAME).exe main.go

.PHONY: build-linux
build-linux:
	GOOS=linux go build -mod=readonly -v \
		--ldflags="-w -s \
			-X $(MODULE)/cmd.LogLevel=panic \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/release/linux/$(GOARCH)/$(NAME) main.go

.PHONY: build-darwin
build-darwin:
	GOOS=darwin go build -mod=readonly -v \
		--ldflags="-w -s \
			-X $(MODULE)/cmd.LogLevel=panic \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/release/darwin/$(GOARCH)/$(NAME) main.go

.PHONY: build-docker
build-docker:
	docker build --target prod --build-arg VER=$(VERSION) --build-arg REV=$(REVISION) -t $(NAME):$(VERSION) .

.PHONY: build-windows-dev
build-windows-dev:
	GOOS=windows go build -mod=readonly -v \
		--ldflags="-X $(MODULE)/cmd.Debug=true \
			-X $(MODULE)/cmd.LogLevel=debug \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/dev/windows/$(GOARCH)/$(NAME).exe main.go

.PHONY: build-linux-dev
build-linux-dev:
	GOOS=linux go build -mod=readonly -v \
		--ldflags="-X $(MODULE)/cmd.Debug=true \
			-X $(MODULE)/cmd.LogLevel=debug \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/dev/linux/$(GOARCH)/$(NAME) main.go

.PHONY: build-darwin-dev
build-darwin-dev:
	GOOS=darwin go build -mod=readonly -v \
		--ldflags="-X $(MODULE)/cmd.Debug=true \
			-X $(MODULE)/cmd.LogLevel=debug \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/dev/darwin/$(GOARCH)/$(NAME) main.go

.PHONY: build-docker-dev
build-docker-dev:
	docker build --target debug --build-arg VER=$(VERSION) --build-arg REV=$(REVISION) -t $(NAME):$(VERSION)-dev .

.PHONY: get-version
get-version:
	echo $(VERSION)

.PHONY: clean
clean:
	git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf

.PHONY: test
test: bootstrap
test:
	go test -v -mod=readonly -covermode=atomic -coverprofile=coverage.txt \
	-failfast -timeout=2m ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: drone
drone:
	drone fmt --save .drone.yml && \
	drone sign $(REPO) --save

.PHONY: generate
generate:
	go generate -x ./...
generate: format

.PHONY: format
format:
	gofmt -w .

.PHONY: lint
lint:
	(which golangci-lint || go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.22.2)
	golangci-lint run ./...
