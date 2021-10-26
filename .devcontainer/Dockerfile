FROM ubuntu:21.04

ARG BAZELISK_SHA256_HASH="4cb534c52cdd47a6223d4596d530e7c9c785438ab3b0a49ff347e991c210b2cd"
ARG BAZELISK_VERSION="v1.10.1"

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -q -y \
        python \
        python3 \
        wget \
        ca-certificates \
        gcc \
        git \
    && apt-get clean \
    && wget -q -O /tmp/bazelisk https://github.com/bazelbuild/bazelisk/releases/download/${BAZELISK_VERSION}/bazelisk-linux-amd64 \
    && echo "${BAZELISK_SHA256_HASH}  /tmp/bazelisk" | sha256sum -c - \
    && mv /tmp/bazelisk /usr/local/bin/bazel \
    && chmod +x /usr/local/bin/bazel
