<a name="unreleased"></a>
## [Unreleased]

### Chore
- refactoring dev tools
- **deps:** update dependency io_bazel_rules_go to v0.24.1
- **deps:** updating AWS session-manager-plugin to version 1.2.7.0
- **deps:** update dependency bazel_gazelle to v0.22.0
- **deps:** update dependency io_bazel_rules_go to v0.24.5
- **deps:** update module spf13/cobra to v1.1.1
- **deps:** update module spf13/cobra to v1.1.0
- **deps:** update module sirupsen/logrus to v1.7.0
- **deps:** bumping go version to 1.15.3
- **deps:** update dependency io_bazel_rules_go to v0.24.4
- **deps:** update dependency io_bazel_rules_docker to v0.15.0
- **deps:** update dependency bazel_gazelle to v0.22.2
- **deps:** update module aws/aws-sdk-go to v1.34.32
- **deps:** update dependency io_bazel_rules_go to v0.24.3
- **deps:** update dependency bazel_gazelle to v0.22.1
- **deps:** update dependency io_bazel_rules_go to v0.24.2
- **deps:** update l.gcr.io/google/bazel docker tag to v3.5.0
- **renovate:** scheduled checks
- **tools:** updating dev tools

### Fix
- **bazel:** semi-hermitizing Go SDK
- **drone:** updating drone.yml signature

### Test
- **aws:** adding unittests

### Update
- **go:** version 1.15.2
- **go:** version 1.15.1


<a name="0.7.0"></a>
## [0.7.0] - 2020-08-27
### Chore
- **release:** version 0.7.0

### Feat
- **cmd:** moving target to arg instead of separate flag

### Fix
- **aws:** error log on a send ssh pub key failure
- **cmd:** clarify ssh usage
- **docs:** adding profile config entry description


<a name="0.6.1"></a>
## [0.6.1] - 2020-08-19
### Chore
- **release:** version 0.6.1

### Fix
- **aws:** reducing timeouts to speed up error feedback loop


<a name="0.6.0"></a>
## [0.6.0] - 2020-08-16
### Aws
- customizing default retryer for AWS api calls

### Chore
- removing deprecated name tag reference
- **deps:** update golang.org/x/crypto commit hash to 75b2880
- **deps:** update Go to 1.15
- **deps:** bumping Go and tidying dependecies
- **deps:** update module aws/aws-sdk-go to v1.34.2
- **deps:** update dependency io_bazel_rules_go to v0.23.7
- **deps:** update module spf13/viper to v1.7.1
- **deps:** update module golang/mock to v1.4.4
- **deps:** update golang.org/x/crypto commit hash to 123391f
- **deps:** update l.gcr.io/google/bazel docker tag to v3.4.1
- **deps:** update dependency io_bazel_rules_docker to v0.14.4
- **deps:** x/crypto version bump
- **release:** version 0.6.0

### Fix
- **bazel:** docker rules
- **pre-commit:** adding Gazelle run during fmt phase


<a name="0.5.3"></a>
## [0.5.3] - 2020-06-21
### Chore
- **bazel:** bump to version 3.3.0
- **deps:** update go modules in Bazel
- **deps:** update module aws/aws-sdk-go to v1.32.6
- **deps:** update l.gcr.io/google/bazel docker tag to v3.3.0
- **pre-commit:** add go mod tidy to update_deps
- **release:** version 0.5.3

### Rollback
- name-tag target type change

### Update
- **go:** version 1.14.4


<a name="0.5.2"></a>
## [0.5.2] - 2020-06-06
### Build
- **deps:** bump github.com/aws/aws-sdk-go from 1.31.7 to 1.31.11
- **deps:** bump github.com/aws/aws-sdk-go from 1.31.1 to 1.31.7
- **deps:** bump github.com/stretchr/testify from 1.5.1 to 1.6.0

### Chore
- **deps:** add renovate.json
- **deps:** update module stretchr/testify to v1.6.1
- **deps:** update dependency io_bazel_rules_go to v0.23.3
- **deps:** update dependency io_bazel_rules_docker to v0.14.3
- **deps:** update dependency bazel_gazelle to v0.21.1
- **deps:** update golang.org/x/crypto commit hash to 70a84ac
- **release:** 0.5.2

### Fix
- **ssh:** error handling

### Update
- **bazel:** to version 3.2.0
- **deps:** Go dependecies in bazel


<a name="0.5.1"></a>
## [0.5.1] - 2020-05-23
### Add
- **aws:** append AWS UA with sigil version

### Build
- **deps:** bump github.com/aws/aws-sdk-go from 1.31.0 to 1.31.1
- **deps:** bump github.com/aws/aws-sdk-go from 1.30.29 to 1.31.0
- **deps:** bump github.com/aws/aws-sdk-go from 1.30.28 to 1.30.29
- **deps:** bump github.com/aws/aws-sdk-go from 1.30.27 to 1.30.28
- **deps:** bump github.com/aws/aws-sdk-go from 1.30.26 to 1.30.27
- **deps:** bump github.com/aws/aws-sdk-go from 1.30.25 to 1.30.26
- **deps:** bump gopkg.in/yaml.v2 from 2.2.8 to 2.3.0
- **deps:** bump github.com/aws/aws-sdk-go from 1.30.24 to 1.30.25
- **deps:** bump github.com/spf13/viper from 1.6.3 to 1.7.0
- **deps:** bump github.com/aws/aws-sdk-go from 1.30.9 to 1.30.24
- **deps:** bump github.com/sirupsen/logrus from 1.5.0 to 1.6.0

### Chore
- using Bazel as a build system
- **release:** 0.5.1

### Fix
- app exit on failed session termination
- version in makefile
- **drone:** docker-release step
- **lint:** addressing comments from deepsource.io

### Rm
- **changelog:** merge PR commits
- **drone:** release notes

### Update
- dependencies


<a name="0.5.0"></a>
## [0.5.0] - 2020-04-18
### Chore
- **pre-commit:** always run

### Delete
- **doc:** manual pages

### Feat
- adding support for environment variables
- **ssh:** allowing custom temp. key directories

### Fix
- **doc:** adding missing flag in ssh_config example
- **lint:** removing broken bin path


<a name="0.4.1"></a>
## [0.4.1] - 2020-04-16
### Fix
- **list:** filters

### Update
- **version:** 0.4.1


<a name="0.4.0"></a>
## [0.4.0] - 2020-04-11
### Add
- deepsource integration
- **desc:** for tests
- **pkg:** unit tests
- **pre-commit:** adding pre-commit and dev section

### Feat
- **golangci-lint:** added config file

### Fix
- **ssh:** adding missing mfa token

### Update
- **doc:** expanding documentation
- **go:** bumping go and dependencies versions
- **version:** to 0.4.0


<a name="0.3.3"></a>
## [0.3.3] - 2020-03-19

<a name="0.3.2"></a>
## [0.3.2] - 2020-03-11

<a name="0.3.1"></a>
## [0.3.1] - 2019-07-18

<a name="0.3.0"></a>
## [0.3.0] - 2019-07-13
### Stargate
- Adding Support for SSH and SCP ([#44](https://github.com/danmx/sigil/issues/44))


<a name="0.2.1"></a>
## [0.2.1] - 2019-05-14

<a name="0.2.0"></a>
## [0.2.0] - 2019-05-03

<a name="0.1.2"></a>
## [0.1.2] - 2019-04-29

<a name="0.1.1"></a>
## [0.1.1] - 2019-04-23

<a name="0.1.0"></a>
## [0.1.0] - 2019-04-23

<a name="0.0.8"></a>
## [0.0.8] - 2019-04-22

<a name="0.0.7"></a>
## [0.0.7] - 2019-04-16

<a name="0.0.6"></a>
## [0.0.6] - 2019-04-16

<a name="0.0.5"></a>
## [0.0.5] - 2019-04-15

<a name="0.0.4"></a>
## [0.0.4] - 2019-04-15

<a name="0.0.3"></a>
## [0.0.3] - 2019-03-19

<a name="0.0.2"></a>
## [0.0.2] - 2019-03-19

<a name="0.0.1"></a>
## 0.0.1 - 2019-03-18

[Unreleased]: https://github.com/danmx/sigil/compare/0.7.0...HEAD
[0.7.0]: https://github.com/danmx/sigil/compare/0.6.1...0.7.0
[0.6.1]: https://github.com/danmx/sigil/compare/0.6.0...0.6.1
[0.6.0]: https://github.com/danmx/sigil/compare/0.5.3...0.6.0
[0.5.3]: https://github.com/danmx/sigil/compare/0.5.2...0.5.3
[0.5.2]: https://github.com/danmx/sigil/compare/0.5.1...0.5.2
[0.5.1]: https://github.com/danmx/sigil/compare/0.5.0...0.5.1
[0.5.0]: https://github.com/danmx/sigil/compare/0.4.1...0.5.0
[0.4.1]: https://github.com/danmx/sigil/compare/0.4.0...0.4.1
[0.4.0]: https://github.com/danmx/sigil/compare/0.3.3...0.4.0
[0.3.3]: https://github.com/danmx/sigil/compare/0.3.2...0.3.3
[0.3.2]: https://github.com/danmx/sigil/compare/0.3.1...0.3.2
[0.3.1]: https://github.com/danmx/sigil/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/danmx/sigil/compare/0.2.1...0.3.0
[0.2.1]: https://github.com/danmx/sigil/compare/0.2.0...0.2.1
[0.2.0]: https://github.com/danmx/sigil/compare/0.1.2...0.2.0
[0.1.2]: https://github.com/danmx/sigil/compare/0.1.1...0.1.2
[0.1.1]: https://github.com/danmx/sigil/compare/0.1.0...0.1.1
[0.1.0]: https://github.com/danmx/sigil/compare/0.0.8...0.1.0
[0.0.8]: https://github.com/danmx/sigil/compare/0.0.7...0.0.8
[0.0.7]: https://github.com/danmx/sigil/compare/0.0.6...0.0.7
[0.0.6]: https://github.com/danmx/sigil/compare/0.0.5...0.0.6
[0.0.5]: https://github.com/danmx/sigil/compare/0.0.4...0.0.5
[0.0.4]: https://github.com/danmx/sigil/compare/0.0.3...0.0.4
[0.0.3]: https://github.com/danmx/sigil/compare/0.0.2...0.0.3
[0.0.2]: https://github.com/danmx/sigil/compare/0.0.1...0.0.2
