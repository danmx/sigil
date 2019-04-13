# sigil
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdanmx%2Fsigil.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdanmx%2Fsigil?ref=badge_shield)
[![Build Status](https://cloud.drone.io/api/badges/danmx/sigil/status.svg)](https://cloud.drone.io/danmx/sigil)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b4725f567cbf46a493a5436ee698b571)](https://www.codacy.com/app/danmx/sigil?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=danmx/sigil&amp;utm_campaign=Badge_Grade) 
[![codecov](https://codecov.io/gh/danmx/sigil/branch/master/graph/badge.svg)](https://codecov.io/gh/danmx/sigil)

## Description

> *Sigil* is the hub of the Great Wheel, a city at the center of the Outlands, the most balanced of neutral areas at the center of the planes. Also known as the "City of Doors" for the multitude of portals to other planes of existence and the Cage since those portals are the only way in or out, it is the setting for most of Planescape: Torment.

*Sigil* is an AWS SSM Session manager client inspired by [xen0l's aws-gate](https://github.com/xen0l/aws-gate).

## Features

- configuration files (YAML) support
- lightweight [container image](https://hub.docker.com/r/danmx/sigil) (~16MB)

## External dependencies

- AWS [session-manager-plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)
- target EC2 instance must have AWS SSM Agent installed ([full guide](https://docs.aws.amazon.com/systems-manager/latest/userguide/ssm-agent.html))
- target EC2's instance profile should have **AmazonEC2RoleforSSM** managed IAM policy attached
  or other with similar permissions ([Adding Session Manager Permissions to an Existing Instance Profile](https://docs.aws.amazon.com/systems-manager/latest/userguide/getting-started-add-permissions-to-existing-profile.html))

## Examples

### Usage

Docker:

```console
docker run --rm -it -v "${HOME}"/.sigil:/home/.sigil -v "${HOME}"/.aws:/home/.aws danmx/sigil:0.0 list --output-format wide
```

Binary:

```console
sigil -r eu-west-1 session --type instance-id --target i-xxxxxxxxxxxxxxxxx
```

Using with [aws-vault](https://github.com/99designs/aws-vault):

```console
aws-vault exec AWS_PROFILE -- sigil -r eu-west-1 session --type instance-id --target i-xxxxxxxxxxxxxxxxx
```

### Config file

By default configuration files are stored in `$HOME/.sigil/` and the default config file is `config.yaml`.

```yaml
region: eu-west-1
target: worker-Node
type: name-tag                  # instance-id/priv-dns/name-tag
output-format: wide             # text/json/yaml/wide
```

## Build

### Binaries

To build binaries for all platforms (Linux, Mac, Windows) and Docker image run:

```console
make build
```

To run specific build use:

```console
make build-[linux|mac|windows]
```

Binaries are located in:

- Linux: `bin/release/linux/amd64/sigil`
- Mac: `bin/release/darwin/amd64/sigil`
- Windows: `bin/release/darwin/amd64/sigil.exe`

### Docker

To only build docker image run:

```console
make build-docker
```

It'll create a docker image tagged `sigil:{version}` where `{version}` corresponds to sigil's current version.

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdanmx%2Fsigil.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdanmx%2Fsigil?ref=badge_large)
