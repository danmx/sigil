<a name="unreleased"></a>
## [Unreleased]


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
- **lint:** addressing comments from deepsource.io

### Update
- dependencies

### Pull Requests
- Merge pull request [#87](https://github.com/danmx/sigil/issues/87) from danmx/bazel
- Merge pull request [#84](https://github.com/danmx/sigil/issues/84) from danmx/dependabot/go_modules/github.com/aws/aws-sdk-go-1.31.1
- Merge pull request [#83](https://github.com/danmx/sigil/issues/83) from danmx/dependabot/go_modules/github.com/aws/aws-sdk-go-1.31.0
- Merge pull request [#82](https://github.com/danmx/sigil/issues/82) from danmx/dependabot/go_modules/github.com/aws/aws-sdk-go-1.30.29
- Merge pull request [#81](https://github.com/danmx/sigil/issues/81) from danmx/dependabot/go_modules/github.com/aws/aws-sdk-go-1.30.28
- Merge pull request [#80](https://github.com/danmx/sigil/issues/80) from danmx/dependabot/go_modules/github.com/aws/aws-sdk-go-1.30.27
- Merge pull request [#79](https://github.com/danmx/sigil/issues/79) from danmx/dependabot/go_modules/github.com/aws/aws-sdk-go-1.30.26
- Merge pull request [#78](https://github.com/danmx/sigil/issues/78) from danmx/dependabot/go_modules/gopkg.in/yaml.v2-2.3.0
- Merge pull request [#77](https://github.com/danmx/sigil/issues/77) from danmx/dependabot/go_modules/github.com/aws/aws-sdk-go-1.30.25
- Merge pull request [#76](https://github.com/danmx/sigil/issues/76) from danmx/dependabot/go_modules/github.com/spf13/viper-1.7.0
- Merge pull request [#73](https://github.com/danmx/sigil/issues/73) from danmx/dependabot/go_modules/github.com/aws/aws-sdk-go-1.30.24
- Merge pull request [#74](https://github.com/danmx/sigil/issues/74) from danmx/dependabot/go_modules/github.com/sirupsen/logrus-1.6.0
- Merge pull request [#72](https://github.com/danmx/sigil/issues/72) from danmx/log
- Merge pull request [#71](https://github.com/danmx/sigil/issues/71) from danmx/append-aws-ua
- Merge pull request [#70](https://github.com/danmx/sigil/issues/70) from danmx/bump-dependencies


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

### Pull Requests
- Merge pull request [#69](https://github.com/danmx/sigil/issues/69) from danmx/hotfix-readme
- Merge pull request [#68](https://github.com/danmx/sigil/issues/68) from danmx/bug-[#65](https://github.com/danmx/sigil/issues/65)
- Merge pull request [#67](https://github.com/danmx/sigil/issues/67) from danmx/bug-63
- Merge pull request [#66](https://github.com/danmx/sigil/issues/66) from danmx/ssh-doc


<a name="0.4.1"></a>
## [0.4.1] - 2020-04-16
### Fix
- **list:** filters

### Update
- **version:** 0.4.1

### Pull Requests
- Merge pull request [#62](https://github.com/danmx/sigil/issues/62) from danmx/tag-fix


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

### Pull Requests
- Merge pull request [#60](https://github.com/danmx/sigil/issues/60) from danmx/v0.4.0
- Merge pull request [#59](https://github.com/danmx/sigil/issues/59) from danmx/refactor
- Merge pull request [#58](https://github.com/danmx/sigil/issues/58) from ekesken/patch-1


<a name="0.3.3"></a>
## [0.3.3] - 2020-03-19
### Pull Requests
- Merge pull request [#57](https://github.com/danmx/sigil/issues/57) from efpe/fixctrlz


<a name="0.3.2"></a>
## [0.3.2] - 2020-03-11
### Pull Requests
- Merge pull request [#56](https://github.com/danmx/sigil/issues/56) from zscholl/zscholl/release-0.3.2
- Merge pull request [#55](https://github.com/danmx/sigil/issues/55) from zscholl/zscholl/pin-version
- Merge pull request [#54](https://github.com/danmx/sigil/issues/54) from zscholl/zscholl/filter-running-instances
- Merge pull request [#53](https://github.com/danmx/sigil/issues/53) from zscholl/master


<a name="0.3.1"></a>
## [0.3.1] - 2019-07-18

<a name="0.3.0"></a>
## [0.3.0] - 2019-07-13
### Stargate
- Adding Support for SSH and SCP ([#44](https://github.com/danmx/sigil/issues/44))


<a name="0.2.1"></a>
## [0.2.1] - 2019-05-14
### Pull Requests
- Merge pull request [#43](https://github.com/danmx/sigil/issues/43) from danmx/hotfix


<a name="0.2.0"></a>
## [0.2.0] - 2019-05-03
### Pull Requests
- Merge pull request [#42](https://github.com/danmx/sigil/issues/42) from danmx/sessions


<a name="0.1.2"></a>
## [0.1.2] - 2019-04-29
### Pull Requests
- Merge pull request [#38](https://github.com/danmx/sigil/issues/38) from danmx/homebrew


<a name="0.1.1"></a>
## [0.1.1] - 2019-04-23
### Pull Requests
- Merge pull request [#33](https://github.com/danmx/sigil/issues/33) from danmx/fix-release
- Merge pull request [#32](https://github.com/danmx/sigil/issues/32) from danmx/optional-config
- Merge pull request [#31](https://github.com/danmx/sigil/issues/31) from danmx/test-release


<a name="0.1.0"></a>
## [0.1.0] - 2019-04-23
### Pull Requests
- Merge pull request [#30](https://github.com/danmx/sigil/issues/30) from danmx/rollback-zip
- Merge pull request [#29](https://github.com/danmx/sigil/issues/29) from danmx/cli-refactor


<a name="0.0.8"></a>
## [0.0.8] - 2019-04-22
### Pull Requests
- Merge pull request [#28](https://github.com/danmx/sigil/issues/28) from danmx/mfa


<a name="0.0.7"></a>
## [0.0.7] - 2019-04-16
### Pull Requests
- Merge pull request [#26](https://github.com/danmx/sigil/issues/26) from danmx/profile-region


<a name="0.0.6"></a>
## [0.0.6] - 2019-04-16
### Pull Requests
- Merge pull request [#25](https://github.com/danmx/sigil/issues/25) from danmx/fix


<a name="0.0.5"></a>
## [0.0.5] - 2019-04-15
### Pull Requests
- Merge pull request [#24](https://github.com/danmx/sigil/issues/24) from danmx/profile
- Merge pull request [#23](https://github.com/danmx/sigil/issues/23) from danmx/log-level


<a name="0.0.4"></a>
## [0.0.4] - 2019-04-15
### Pull Requests
- Merge pull request [#22](https://github.com/danmx/sigil/issues/22) from danmx/fix-docker


<a name="0.0.3"></a>
## [0.0.3] - 2019-03-19

<a name="0.0.2"></a>
## [0.0.2] - 2019-03-19

<a name="0.0.1"></a>
## 0.0.1 - 2019-03-18

[Unreleased]: https://github.com/danmx/sigil/compare/0.5.1...HEAD
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
