#!/bin/sh

set -eu

test(){
    bazel test //...
}

fmt() {
    bazel run :gazelle && bazel run @go_sdk//:bin/gofmt -- -s -w .
}

lint() {
    bazel run @go_sdk//:bin/go -- generate -x ./... && bazel run @com_github_danmx_bazel_tools//golangci-lint:run -- run ./...
}

update_deps() {
    bazel run @go_sdk//:bin/go -- mod tidy && bazelisk run :gazelle -- update-repos -from_file=go.mod -to_macro=tools/repositories.bzl%go_repositories -prune=true
}

ci() {
    bazel run @com_github_danmx_bazel_tools//drone-cli:run -- fmt --save .drone.yml && bazel run @com_github_danmx_bazel_tools//drone-cli:run -- sign danmx/sigil --save
}

changelog(){
    bazel run @com_github_danmx_bazel_tools//git-chglog:run -- -o CHANGELOG.md
}

#shellcheck disable=SC2068
$@
