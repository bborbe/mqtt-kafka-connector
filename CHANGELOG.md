# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.2.1

- Add GitHub Actions CI workflow

## v1.2.0
- Modernize build tooling and linting infrastructure
- Replace deprecated golint with golangci-lint
- Add comprehensive security scanning (gosec, osv-scanner, trivy)
- Add code formatting tools (golines, go-modtool)
- Migrate to github.com/bborbe/http for server with secure defaults
- Update to github.com/bborbe/argument/v2 with context support
- Update to github.com/bborbe/errors with context support
- Remove vendor directory from git tracking
- Add .gitignore for vendor directory
- Update Makefile to use -mod=mod instead of -mod=vendor

## v1.1.9
- Update dependencies to latest versions
- Update bborbe/run and bborbe/errors libraries
- Update Kafka client (IBM/sarama)
- Update Ginkgo and Gomega testing frameworks
- Remove deprecated dependencies (raven-go, automaxprocs, certifi)

## v1.1.8

- go mod update

## v1.1.7

- go mod update

## v1.1.6

- go mod update

## v1.1.5

- go mod update

## v1.1.4

- go mod update

## v1.1.3

- go mod update

## v1.1.2

- go mod update

## v1.1.1

- add vulncheck
- go mod update

## v1.1.0

- go mod update

## v1.0.0

- Initial version
