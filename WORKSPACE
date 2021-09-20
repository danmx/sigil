# gazelle:repository_macro tools/repositories.bzl%go_repositories

workspace(
    name = "com_github_danmx_sigil",
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

# Golang
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "8e968b5fcea1d2d64071872b12737bbb5514524ee5f0a4f54f5920266c261acb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.28.0/rules_go-v0.28.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.28.0/rules_go-v0.28.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_download_sdk", "go_register_toolchains", "go_rules_dependencies")

go_download_sdk(
    name = "go_sdk",
    version = "1.15.3",
)

go_rules_dependencies()

go_register_toolchains()

# gazelle
http_archive(
    name = "bazel_gazelle",
    sha256 = "62ca106be173579c0a167deb23358fdfe71ffa1e4cfdddf5582af26520f1c66f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# Container image
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "1f4e59843b61981a96835dc4ac377ad4da9f8c334ebe5e0bb3f58f80c09735f4",
    strip_prefix = "rules_docker-0.19.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.19.0/rules_docker-v0.19.0.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

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
    sha256 = "82a5fb946d2eb0fed80d3d70c2556784ec6cb5c35cd65a1b5e93e46f99681650",
    strip_prefix = "bazel_gomock-1.3",
    urls = ["https://github.com/jmhodges/bazel_gomock/archive/v1.3.tar.gz"],
)

# AWS Session Manager Plugin
http_file(
    name = "session_manager_plugin_deb",
    downloaded_file_path = "session-manager-plugin.deb",
    sha256 = "f1c03d2aaad9f89f73fc70f1c1cdef0e2877a03b86cca3c8b5c97992c6344449",
    urls = ["https://s3.amazonaws.com/session-manager-downloads/plugin/1.2.245.0/ubuntu_64bit/session-manager-plugin.deb"],
)

# golangci-lint, drone-cli & git-chglog
http_archive(
    name = "com_github_danmx_bazel_tools",
    sha256 = "2a21bc87e5b8668b401761b835047f03ed617c9a0398e88eaf3883b6596ff6ed",
    strip_prefix = "bazel-tools-0.1.0",
    urls = ["https://github.com/danmx/bazel-tools/archive/0.1.0.tar.gz"],
)

load("@com_github_danmx_bazel_tools//git-chglog:deps.bzl", "git_chglog_dependencies")
load("@com_github_danmx_bazel_tools//drone-cli:deps.bzl", "drone_cli_dependencies")
load("@com_github_danmx_bazel_tools//golangci-lint:deps.bzl", "golangci_lint_dependencies")

git_chglog_dependencies()

drone_cli_dependencies()

golangci_lint_dependencies()
