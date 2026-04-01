# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.2.13

- Fix golangci-lint compilation by pinning go-header to v0.5.0

## v1.2.12

- Update dependencies to fix security vulnerabilities (go-git/v5 v5.17.2, buildkit v0.29.0)

## v1.2.11

- Update go-git/go-git to v5.17.1 (fix security vulnerabilities)

## v1.2.10

- Bump bborbe/http to v1.26.8 and bborbe/run to v1.9.12
- Update multiple indirect dependencies (otel, prometheus, google APIs, moby)
- Add opencontainers/runtime-spec replace directive
- Enable --allow-parallel-runners for golangci-lint

## v1.2.9

- Pin runtime-spec to v1.2.1 to fix containerd v1.7.30 compile error in CI

## v1.2.8

- Update dependencies to fix containerd v1.7.30 compile error in osv-scanner
- Upgrade osv-scanner v2.3.5, golangci-lint v2.11.4, go-modtool v0.7.1
- Upgrade bborbe libs (argument v2.12.12, errors v1.5.8, http v1.26.7, run v1.9.11)

## v1.2.7

- remove vendored containerd patches directory and local replace directive
- standardize Makefile: multiline trivy, add .PHONY declarations

## v1.2.6

- chore: Verify project health — all tests pass, linting succeeds, and precommit checks are clean

## v1.2.5

- fix: Add local patch for containerd v1.7.30 to fix LinuxPids.Limit type incompatibility with opencontainers/runtime-spec v1.3.0
- chore: Update Makefile to exclude patches/ directory from format, addlicense, osv-scanner, and trivy scans

## v1.2.4

- upgrade golangci-lint from v1 to v2
- add trivy ghcr.io db-repository
- update bborbe deps (argument, errors, http, run)

## v1.2.3

- Update go to 1.26.1
- Bump bborbe/argument, errors, run to latest patch versions
- Update golang.org/x/oauth2, sync, sys, crypto dependencies
- Upgrade grpc v1.79.3, otel v1.39.0, osv-scanner v2.3.4
- Remove large exclude/replace blocks and cleanup go.mod

## v1.2.2

- Fix gosec G118: use signal.NotifyContext instead of manual context cancellation
- Update docker/cli to v29.2.0 (GHSA-p436-gjf2-799p)
- Update dependencies to latest versions

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
