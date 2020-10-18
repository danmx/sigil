# gazelle:repository_macro tools/repositories.bzl%go_repositories

workspace(
    name = "com_github_danmx_sigil",
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

# Golang
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "ac03931e56c3b229c145f1a8b2a2ad3e8d8f1af57e43ef28a26123362a1e3c7e",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.24.4/rules_go-v0.24.4.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.24.4/rules_go-v0.24.4.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_download_sdk", "go_register_toolchains", "go_rules_dependencies")

go_download_sdk(
    name = "go_sdk",
    version = "1.15.2",
)

go_rules_dependencies()

go_register_toolchains()

# gazelle
http_archive(
    name = "bazel_gazelle",
    sha256 = "b85f48fa105c4403326e9525ad2b2cc437babaa6e15a3fc0b1dbab0ab064bc7c",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.2/bazel-gazelle-v0.22.2.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.2/bazel-gazelle-v0.22.2.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# Container image
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "4521794f0fba2e20f3bf15846ab5e01d5332e587e9ce81629c7f96c793bb7036",
    strip_prefix = "rules_docker-0.14.4",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.14.4/rules_docker-v0.14.4.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load("@io_bazel_rules_docker//repositories:pip_repositories.bzl", "pip_deps")

pip_deps()

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

# golangci-lint, drone-cli & git-chglog
http_archive(
    name = "com_github_danmx_bazel_tools",
    sha256 = "08e0311293897958438a838777d1293c69212c1fe4ba042bb72362223b72df8c",
    strip_prefix = "bazel-tools-cf03e739f3bb7338189a43e7268381483992990b",
    urls = ["https://github.com/danmx/bazel-tools/archive/cf03e739f3bb7338189a43e7268381483992990b.tar.gz"],
)

load("@com_github_danmx_bazel_tools//git-chglog:deps.bzl", "git_chglog_dependencies")
load("@com_github_danmx_bazel_tools//drone-cli:deps.bzl", "drone_cli_dependencies")
load("@com_github_danmx_bazel_tools//golangci-lint:deps.bzl", "golangci_lint_dependencies")

git_chglog_dependencies()

drone_cli_dependencies()

golangci_lint_dependencies()
