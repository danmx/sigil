# sigil

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdanmx%2Fsigil.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdanmx%2Fsigil?ref=badge_shield)
[![Build Status](https://cloud.drone.io/api/badges/danmx/sigil/status.svg)](https://cloud.drone.io/danmx/sigil)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b4725f567cbf46a493a5436ee698b571)](https://www.codacy.com/app/danmx/sigil?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=danmx/sigil&amp;utm_campaign=Badge_Grade)
[![codecov](https://codecov.io/gh/danmx/sigil/branch/master/graph/badge.svg)](https://codecov.io/gh/danmx/sigil)
[![DeepSource](https://static.deepsource.io/deepsource-badge-light-mini.svg)](https://deepsource.io/gh/danmx/sigil/?ref=repository-badge)

## Description

> *Sigil* is the hub of the Great Wheel, a city at the center of the Outlands, the most balanced of neutral areas at the center of the planes. Also known as the "City of Doors" for the multitude of portals to other planes of existence and the Cage since those portals are the only way in or out, it is the setting for most of Planescape: Torment.

*Sigil* is an AWS SSM Session manager client. Allowing access to EC2 instances without exposing any ports.

## Features

- configuration files support (TOML, YAML, JSON, etc.)
- support for different configuration profiles
- lightweight [container image](https://hub.docker.com/r/danmx/sigil)
- SSH and SCP support

## External dependencies

### Local

- AWS [session-manager-plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html) (version 1.1.17.0+ for SSH support)

### Remote

- target EC2 instance must have AWS SSM Agent installed ([full guide](https://docs.aws.amazon.com/systems-manager/latest/userguide/ssm-agent.html)) (version 2.3.672.0+ for SSH support)
- AWS [ec2-instance-connect](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-connect-set-up.html) to use SSH with your own and/or temporary keys
- to support AWS SSM target EC2 instance profile should have **AmazonSSMManagedInstanceCore** managed IAM policy attached or a specific policy with similar permissions (check [About Policies for a Systems Manager Instance Profile](https://docs.aws.amazon.com/systems-manager/latest/userguide/setup-instance-profile.html) and [About Minimum S3 Bucket Permissions for SSM Agent](https://docs.aws.amazon.com/systems-manager/latest/userguide/ssm-agent-minimum-s3-permissions.html))

## Documentation

The manual can be found [here](docs/README.md).

## Installation

### Homebrew

```shell
brew tap danmx/sigil
brew install sigil
```

or

```shell
brew install danmx/sigil/sigil
```

### Docker

```shell
docker pull danmx/sigil:0.4
```

## Examples

### Usage

Docker:

```shell
docker run --rm -it -v "${HOME}"/.sigil:/home/nonroot/.sigil -v "${HOME}"/.aws:/home/.aws danmx/sigil:0.4 list --output-format wide
```

Binary:

```shell
sigil -r eu-west-1 session --type instance-id --target i-xxxxxxxxxxxxxxxxx
```

Using with [aws-vault](https://github.com/99designs/aws-vault):

```shell
aws-vault exec AWS_PROFILE -- sigil -r eu-west-1 session --type instance-id --target i-xxxxxxxxxxxxxxxxx
```

### SSH integration

Add an entry to your `ssh_config`:

```ssh_config
Host i-* mi-*
    IdentityFile ~/.sigil/temp_key
    IdentitiesOnly yes
    ProxyCommand sigil ssh --target %h --port %p --pub-key "${HOME}"/.sigil/temp_key.pub --gen-key-pair --os-user %r
Host *.compute.internal
    IdentityFile ~/.sigil/temp_key
    IdentitiesOnly yes
    ProxyCommand sigil ssh --type private-dns --target %h --port %p --pub-key "${HOME}"/.sigil/temp_key.pub --gen-key-pair --os-user %r
```

and run:

```shell
ssh ec2-user@i-123456789
```

or

```shell
ssh ec2-user@ip-10-0-0-5.eu-west-1.compute.internal
```

### Config file

By default configuration file is located in `${HOME}/.sigil/config.toml`.

```toml
[default]
  type = "instance-id"
  output-format = "wide"
  region = "eu-west-1"
  profile = "dev"
  interactive = true
```

## Build

### Binaries

To build binaries for all platforms (Linux, Mac, Windows) and Docker image run:

```shell
make build
```

To run specific build use:

```shell
make build-[linux|darwin|windows]
```

Binaries are located in:

- Linux: `bin/release/linux/amd64/sigil`
- Darwin: `bin/release/darwin/amd64/sigil`
- Windows: `bin/release/windows/amd64/sigil.exe`

### Container image

To only build docker image run:

```shell
make build-docker
```

It'll create a docker image tagged `sigil:{version}` where `{version}` corresponds to sigil's current version.

## Contributions

All contributions are welcomed!

### Dev Dependencies

- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://github.com/golangci/golangci-lint)
- [make](https://www.gnu.org/software/make/)

### Commits

I'm trying to follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

### Bootstraping

```sh
pre-commit install
pre-commit install --hook-type pre-push
make bootstrap
```

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdanmx%2Fsigil.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdanmx%2Fsigil?ref=badge_large)

[Apache 2.0](LICENSE)

## Considerations

*Sigil* was inspired by [xen0l's aws-gate](https://github.com/xen0l/aws-gate).
