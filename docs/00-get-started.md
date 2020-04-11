# Get Started

## External dependencies

To start using `sigil` you need to make sure you have all the necessary dependencies.

### Local

- AWS [session-manager-plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html) (version 1.1.17.0+ for SSH support)

### Remote

- target EC2 instance must have AWS SSM Agent installed ([full guide](https://docs.aws.amazon.com/systems-manager/latest/userguide/ssm-agent.html)) (version 2.3.672.0+ for SSH support)
- AWS [ec2-instance-connect](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-connect-set-up.html) to use SSH with your own and/or temporary keys
- target EC2 instance profile should have **AmazonSSMManagedInstanceCore** managed IAM policy attached or a specific policy with similar permissions (check [About Policies for a Systems Manager Instance Profile](https://docs.aws.amazon.com/systems-manager/latest/userguide/setup-instance-profile.html) and [About Minimum S3 Bucket Permissions for SSM Agent](https://docs.aws.amazon.com/systems-manager/latest/userguide/ssm-agent-minimum-s3-permissions.html))

## Download

To download `sigil` you can use:

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
docker pull danmx/sigil:0.3
```

### Source code

Pull the repository and build binaries.

```shell
git clone https://github.com/danmx/sigil.git
cd sigil
```

For all platforms (Linux, Mac, Windows) and Docker image run:

```shell
make build
```

To run specific build use:

```shell
make build-[linux|darwin|windows|docker]
```

Binaries are located in:

- Linux: `bin/release/linux/amd64/sigil`
- Darwin: `bin/release/darwin/amd64/sigil`
- Windows: `bin/release/darwin/amd64/sigil.exe`
