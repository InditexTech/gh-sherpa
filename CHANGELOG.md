# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.1] - 2025-10-13

### Fixed

- [#132](https://github.com/InditexTech/gh-sherpa/issues/132) Create Pull Request from fork fails when setting labels if you don't have permissions

## [1.4.0] - 2025-07-24

### Added

- [#121](https://github.com/InditexTech/gh-sherpa/issues/121) Fork support for external contributors

## [1.3.1] - 2025-05-29

### Fixed

- [#119](https://github.com/InditexTech/gh-sherpa/issues/119) Prevent branches ending with dash

## [1.3.0] - 2025-04-04

### Added

- [#114](https://github.com/InditexTech/gh-sherpa/issues/114) Add Pull Request template support for create-pr command
- [#113](https://github.com/InditexTech/gh-sherpa/issues/113) Update SECURITY.md

## [1.2.0] - 2024-06-11

### Added

- [#94](https://github.com/InditexTech/gh-sherpa/issues/94) Branch naming not working for types other than feature
- [#93](https://github.com/InditexTech/gh-sherpa/issues/93) Branch Name Length Control Configurable
- [#82](https://github.com/InditexTech/gh-sherpa/issues/82) Add REUSE Compliance workflow

### Refactored

- [#91](https://github.com/InditexTech/gh-sherpa/issues/91) Refactor flow and extract nested functionality
- [#90](https://github.com/InditexTech/gh-sherpa/issues/90) Refactor domain issue and avoid injection
- [#86](https://github.com/InditexTech/gh-sherpa/issues/86) Use fakes instead of mocks for use case tests

### Fixed

- [#95](https://github.com/InditexTech/gh-sherpa/issues/95) Cannot create pull request if remote branch already exists
- [#88](https://github.com/InditexTech/gh-sherpa/issues/88) Pull Request ID as Issue ID
- [#81](https://github.com/InditexTech/gh-sherpa/issues/81) Prevent pushing local branch if remote branch already exists with the same name
- [#79](https://github.com/InditexTech/gh-sherpa/issues/79) Warns about duplicated branches but that's not the case

## [1.1.1] - 2024-02-13

### Documentation

- [#64](https://github.com/InditexTech/gh-sherpa/issues/64) Add `SECURITY.md` file

### Fixed

- [#68](https://github.com/InditexTech/gh-sherpa/issues/68) Version is not shown when gh extension list

## [1.1.0] - 2024-01-16

### Added

- [#58](https://github.com/InditexTech/gh-sherpa/issues/58) Add `kind/internal` label to `internal` issue label in default config
- [#51](https://github.com/InditexTech/gh-sherpa/issues/51) Add GitHub issue type label to the generated pull request

### Documentation

- [#57](https://github.com/InditexTech/gh-sherpa/issues/57) Add CLA signature requirement to contribute in CONTRIBUTING.md file

### Fixed

- [#43](https://github.com/InditexTech/gh-sherpa/issues/43) 404 when accessing configuration link

## [1.0.0] - 2023-12-22

### Added

- [#40](https://github.com/InditexTech/gh-sherpa/issues/40) Add `kind/dependency` label to `dependency` issue label in default config
- [#20](https://github.com/InditexTech/gh-sherpa/issues/20) Validate configuration
- [#18](https://github.com/InditexTech/gh-sherpa/issues/18) Configuration file with comments
- [#9](https://github.com/InditexTech/gh-sherpa/issues/9) Insecure TLS enabled by default

### Dependencies

- [#41](https://github.com/InditexTech/gh-sherpa/issues/41) Bump `github.com/go-playground/validator/v10` to latest version
- [#38](https://github.com/InditexTech/gh-sherpa/issues/38) Bump `github.com/spf13/viper` to latest version (`v1.18.2`)
- [#37](https://github.com/InditexTech/gh-sherpa/issues/37) Bump `github.com/spf13/cobra` to latest version (`v1.8.0`)
- [#30](https://github.com/InditexTech/gh-sherpa/issues/30) New release of go-gh

### Refactored

- [#36](https://github.com/InditexTech/gh-sherpa/issues/36) Update to go1.21.5 and replace `golang.org/x/exp/slices` with the `slices` in the std library
- [#11](https://github.com/InditexTech/gh-sherpa/issues/11) Allow branch prefix with issue type mapping
- [#10](https://github.com/InditexTech/gh-sherpa/issues/10) Extract interactive flag to a global level

### Documentation

- [#15](https://github.com/InditexTech/gh-sherpa/issues/15) Improve and update documentation

[Unreleased]: https://github.com/InditexTech/gh-sherpa/compare/1.3.1...main
[1.3.1]: https://github.com/InditexTech/gh-sherpa/compare/1.3.0...1.3.1
[1.3.0]: https://github.com/InditexTech/gh-sherpa/compare/1.2.0...1.3.0
[1.2.0]: https://github.com/InditexTech/gh-sherpa/compare/1.1.1...1.2.0
[1.1.1]: https://github.com/InditexTech/gh-sherpa/compare/1.1.0...1.1.1
[1.1.0]: https://github.com/InditexTech/gh-sherpa/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/InditexTech/gh-sherpa/commits/1.0.0
