
# Change this and commit to create new release
override VERSION ?= 0.2.1
NAME = sigil
REPO = danmx/$(NAME)
MODULE = github.com/$(REPO)
override REVISION ?= $(shell git rev-parse HEAD;)

export GO111MODULE = on

.PHONY: bootstrap
bootstrap:
	@go mod download && go mod vendor

.PHONY: build
build: bootstrap build-linux build-mac build-windows

.PHONY: build-dev
build-dev: bootstrap build-linux-dev build-mac-dev build-windows-dev

.PHONY: release
release: build release-windows release-linux release-darwin

release-windows:
	@mkdir -p dist && tar -czvf dist/$(NAME)_windows-amd64.tar.gz -C bin/release/windows/amd64/ $(NAME).exe

release-linux:
	@mkdir -p dist && tar -czvf dist/$(NAME)_linux-amd64.tar.gz -C bin/release/linux/amd64/ $(NAME)

release-darwin:
	@mkdir -p dist && tar -czvf dist/$(NAME)_darwin-amd64.tar.gz -C bin/release/darwin/amd64/ $(NAME)

build-windows: export GOARCH=amd64
build-windows:
	@GOOS=windows go build -mod=vendor -v \
		--ldflags="-w -s \
			-X $(MODULE)/cmd.LogLevel=panic \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/release/windows/amd64/$(NAME).exe main.go

build-linux: export GOARCH=amd64
build-linux: export CGO_ENABLED=0
build-linux:
	@GOOS=linux go build -mod=vendor -v \
		--ldflags="-w -s \
			-X $(MODULE)/cmd.LogLevel=panic \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/release/linux/amd64/$(NAME) main.go

build-mac: export GOARCH=amd64
build-mac: export CGO_ENABLED=0
build-mac:
	@GOOS=darwin go build -mod=vendor -v \
		--ldflags="-w -s \
			-X $(MODULE)/cmd.LogLevel=panic \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/release/darwin/amd64/$(NAME) main.go

build-docker:
	@docker build --target prod --build-arg VER=$(VERSION) --build-arg REV=$(REVISION) -t $(NAME):$(VERSION) .

build-windows-dev: export GOARCH=amd64
build-windows-dev:
	@GOOS=windows go build -mod=vendor -v \
		--ldflags="-X $(MODULE)/cmd.Debug=true \
			-X $(MODULE)/cmd.LogLevel=debug \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/dev/windows/amd64/$(NAME).exe main.go

build-linux-dev: export GOARCH=amd64
build-linux-dev: export CGO_ENABLED=0
build-linux-dev:
	@GOOS=linux go build -mod=vendor -v \
		--ldflags="-X $(MODULE)/cmd.Debug=true \
			-X $(MODULE)/cmd.LogLevel=debug \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/dev/linux/amd64/$(NAME) main.go

build-mac-dev: export GOARCH=amd64
build-mac-dev: export CGO_ENABLED=0
build-mac-dev:
	@GOOS=darwin go build -mod=vendor -v \
		--ldflags="-X $(MODULE)/cmd.Debug=true \
			-X $(MODULE)/cmd.LogLevel=debug \
			-X $(MODULE)/cmd.AppName=$(NAME) \
			-X $(MODULE)/cmd.Version=$(VERSION) \
			-X $(MODULE)/cmd.Revision=$(REVISION)" \
		-o bin/dev/darwin/amd64/$(NAME) main.go

build-docker-dev:
	@docker build --target debug --build-arg VER=$(VERSION) --build-arg REV=$(REVISION) -t $(NAME):$(VERSION) .

.PHONY: get-version
get-version:
	@echo $(VERSION)

.PHONY: clean
clean:
	@git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf

.PHONY: test
test: bootstrap
test:
	@go test -v -mod=vendor -race -covermode=atomic -coverprofile=coverage.txt \
	-failfast -timeout=2m ./...

.PHONY: tidy
	@go mod tidy

.PHONY: drone-sign
drone-sign:
	@drone fmt --save .drone.yml && drone sign $(REPO) --save
