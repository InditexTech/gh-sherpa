# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/InditexTech/gh-sherpa/compare/1.0.0...develop
[1.0.0]: https://github.com/InditexTech/gh-sherpa/commits/1.0.0
