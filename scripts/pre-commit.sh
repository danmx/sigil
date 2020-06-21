#!/bin/sh

set -eu

test(){
    bazelisk test //...
}

fmt() {
    bazelisk run :gofmt
}

lint() {
    bazelisk run :generate && bazelisk run :lint
}

update_deps() {
    bazelisk run :tidy && bazelisk run :gazelle -- update-repos -from_file=go.mod -to_macro=tools/repositories.bzl%go_repositories -prune=true
}

ci() {
    bazelisk run :drone-fmt && bazelisk run :drone-sign
}

changelog(){
    bazelisk run :changelog
}

#shellcheck disable=SC2068
$@
