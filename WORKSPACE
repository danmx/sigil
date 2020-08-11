# gazelle:repository_macro tools/repositories.bzl%go_repositories

workspace(
    name = "com_github_danmx_sigil",
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

# Golang
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "0310e837aed522875791750de44408ec91046c630374990edd51827cb169f616",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.23.7/rules_go-v0.23.7.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.23.7/rules_go-v0.23.7.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(
    go_version = "1.14.4",
)

# gazelle
http_archive(
    name = "bazel_gazelle",
    sha256 = "cdb02a887a7187ea4d5a27452311a75ed8637379a1287d8eeb952138ea485f7d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.1/bazel-gazelle-v0.21.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.1/bazel-gazelle-v0.21.1.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# Container image
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "6287241e033d247e9da5ff705dd6ef526bac39ae82f3d17de1b69f8cb313f9cd",
    strip_prefix = "rules_docker-0.14.3",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.14.3/rules_docker-v0.14.3.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

load("//tools:repositories.bzl", "go_repositories")

go_repositories()

# GoMock
http_archive(
    name = "com_github_jmhodges_bazel_gomock",
    sha256 = "4baf3389ca48c30d8b072a027923c91c45915ab8061e39e7a0c62706332e096e",
    strip_prefix = "bazel_gomock-1.2",
    urls = ["https://github.com/jmhodges/bazel_gomock/archive/v1.2.tar.gz"],
)

# AWS Session Manager Plugin
http_file(
    name = "session_manager_plugin_deb",
    downloaded_file_path = "session-manager-plugin.deb",
    sha256 = "d4d578a64210165ec434d658212304a968acb2efa49074868552427e738ea97c",
    urls = ["https://s3.amazonaws.com/session-manager-downloads/plugin/1.1.61.0/ubuntu_64bit/session-manager-plugin.deb"],
)

# golangci-lint, drone-cli, git-chglog & workspace_binary
http_archive(
    name = "com_github_danmx_bazel_tools",
    sha256 = "6493f27aba59c1fb91adcc28f267a53b02393a8c7895c9f266c11bd870631c47",
    strip_prefix = "bazel-tools-d35bbdbaecc70c5c033b3af148ee4c8e9b1d31d4",
    urls = ["https://github.com/danmx/bazel-tools/archive/d35bbdbaecc70c5c033b3af148ee4c8e9b1d31d4.tar.gz"],
)

load("@com_github_danmx_bazel_tools//git-chglog:deps.bzl", "git_chglog_dependencies")
load("@com_github_danmx_bazel_tools//drone-cli:deps.bzl", "drone_cli_dependencies")
load("@com_github_danmx_bazel_tools//golangci-lint:deps.bzl", "golangci_lint_dependencies")

git_chglog_dependencies()

drone_cli_dependencies()

golangci_lint_dependencies()
