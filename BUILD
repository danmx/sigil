# gazelle:prefix github.com/danmx/sigil
# gazelle:proto disable_global

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//container:image.bzl", "container_image")
load("@io_bazel_rules_docker//container:layer.bzl", "container_layer")
load("@io_bazel_rules_docker//docker:docker.bzl", "docker_push")
load("@com_github_danmx_bazel_tools//run_in_workspace:def.bzl", "workspace_binary")
load("@com_github_danmx_bazel_tools//git-chglog:def.bzl", "git_chglog")
load("@com_github_danmx_bazel_tools//drone-cli:def.bzl", "drone")
load("@com_github_danmx_bazel_tools//golangci-lint:def.bzl", "golangci_lint")

package(default_visibility = ["//visibility:public"])

golangci_lint(
    name = "lint",
    args = [
        "run",
        "./...",
    ],
)

drone(
    name = "drone-fmt",
    args = [
        "fmt",
        "--save",
        ".drone.yml",
    ],
)

drone(
    name = "drone-sign",
    args = [
        "sign",
        "danmx/sigil",
        "--save",
    ],
)

git_chglog(
    name = "changelog",
    args = [
        "-o",
        "CHANGELOG.md",
    ],
)

genrule(
    name = "concat-cov",
    srcs = glob(["bazel-out/**/testlogs/**/coverage.dat"]),
    outs = ["coverage.txt"],
    cmd_bash = "./$(location //tools:fix_codecov.sh) > \"$@\"",
    exec_tools = ["//tools:fix_codecov.sh"],
)

workspace_binary(
    name = "gofmt",
    args = [
        "-s",
        "-w",
        ".",
    ],
    cmd = "@go_sdk//:bin/gofmt",
)

workspace_binary(
    name = "generate",
    args = [
        "generate",
        "-x",
        "./...",
    ],
    cmd = "@go_sdk//:bin/go",
)

workspace_binary(
    name = "tidy",
    args = [
        "mod",
        "tidy",
    ],
    cmd = "@go_sdk//:bin/go",
)

gazelle(
    name = "gazelle",
    command = "fix",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/danmx/sigil",
    visibility = ["//visibility:private"],
    x_defs = {
        "github.com/danmx/sigil/cmd.gitCommit": "{STABLE_GIT_COMMIT}",
        "github.com/danmx/sigil/cmd.appVersion": "{STABLE_VERSION}",
    },
    deps = ["//cmd:go_default_library"],
)

# Development
go_binary(
    name = "dev",
    out = "dev/sigil",
    embed = [":go_default_library"],
    pure = "on",
    static = "on",
    x_defs = {
        "github.com/danmx/sigil/cmd.gitCommit": "{STABLE_GIT_COMMIT}",
        "github.com/danmx/sigil/cmd.appVersion": "{STABLE_VERSION}",
        "github.com/danmx/sigil/cmd.dev": "true",
        "github.com/danmx/sigil/cmd.logLevel": "debug",
    },
)

# Release
go_binary(
    name = "release",
    out = "sigil",
    embed = [":go_default_library"],
    pure = "on",
    static = "on",
)

pkg_tar(
    name = "sigil_darwin-amd64",
    srcs = [":release"],
    extension = "tar.gz",
    mode = "0o755",
)

pkg_tar(
    name = "sigil_linux-amd64",
    srcs = [":release"],
    extension = "tar.gz",
    mode = "0o755",
)

pkg_tar(
    name = "sigil_windows-amd64",
    srcs = [":release"],
    extension = "zip",
    mode = "0o755",
)

# Include it in our base image as a tar.
container_layer(
    name = "plugin-layer",
    debs = ["@session_manager_plugin_deb//file"],
    symlinks = {"/usr/bin/session-manager-plugin": "/usr/local/sessionmanagerplugin/bin/session-manager-plugin"},
    visibility = ["//visibility:private"],
)

container_layer(
    name = "dev-layer",
    directory = "/usr/bin",
    files = [":dev"],
    visibility = ["//visibility:private"],
)

container_layer(
    name = "release-layer",
    files = [":release"],
    visibility = ["//visibility:private"],
)

container_image(
    name = "dev-image",
    base = "@go_debug_image_base//image",
    cmd = ["--help"],
    entrypoint = ["sigil"],
    layers = [
        "plugin-layer",
        "dev-layer",
    ],
    user = "nonroot",
)

container_image(
    name = "release-image",
    base = "@go_debug_image_base//image",
    cmd = ["--help"],
    entrypoint = ["sigil"],
    layers = [
        "plugin-layer",
        "release-layer",
    ],
    user = "nonroot",
)

docker_push(
    name = "push-dev-image",
    image = ":dev-image",
    registry = "docker.io",
    repository = "danmx/sigil",
    tag = "dev",
)

docker_push(
    name = "push-release-image",
    image = ":release-image",
    registry = "docker.io",
    repository = "danmx/sigil",
    tag = "{STABLE_VERSION}",
)

docker_push(
    name = "push-major-release-image",
    image = ":release-image",
    registry = "docker.io",
    repository = "danmx/sigil",
    tag = "{STABLE_MAJOR_VERSION}",
)

docker_push(
    name = "push-minor-release-image",
    image = ":release-image",
    registry = "docker.io",
    repository = "danmx/sigil",
    tag = "{STABLE_MINOR_VERSION}",
)
