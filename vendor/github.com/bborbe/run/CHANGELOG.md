# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.8.2
- Update Go version from 1.25.2 to 1.25.4 (fixes crypto/x509 performance vulnerability GO-2025-4007)
- Update GitHub Actions CI workflow to use Go 1.25.4
- Update dependencies to latest versions

## v1.8.1
- Add nil validation in FuncRunnerFunc.Run to prevent panics
- Improve receiver name clarity in FuncRunnerFunc (b → f)
- Fix unchecked errors in test files (use proper Gomega assertions)
- Update copyright year in run_background-runner.go (2023-2025)

## v1.8.0
- Add FuncRunner interface for executing functions with custom behavior
- Add FuncRunnerFunc adapter for function-to-interface pattern
- Refactor BackgroundRunner to use FuncRunner interface composition
- Add comprehensive test suite for FuncRunner (8 new tests)
- Update BackgroundRunner documentation to clarify interface embedding
- Update bborbe/errors dependency from v1.3.0 to v1.3.1
- Clean up unused dependencies in go.mod

## v1.7.8
- Update Go version from 1.24.5 to 1.25.2
- Add golangci-lint configuration (.golangci.yml)
- Add security scanning tools (Trivy, gosec, osv-scanner) to Makefile
- Update GitHub Actions workflow with Trivy installation
- Update development dependencies (golangci-lint, osv-scanner, google/addlicense)
- Improve Makefile with new security and linting targets

## v1.7.7

- Code formatting improvements for better readability using golines
- Add golines tool dependency for automated line length management
- Update dependencies and go.mod with latest versions
- Update test files with improved formatting
- Enhance Makefile for better build process

## v1.7.6

- Add comprehensive Go documentation following best practices to all public APIs
- Create package documentation (doc.go) with usage examples 
- Update README with detailed library documentation and examples
- Improve function comments for better godoc rendering
- Update generated mocks (HasCaptureException interface)

## v1.7.5

- go mod update
- update mocks

## v1.7.4

- add tests
- go mod update

## v1.7.3

- add tests

## v1.7.2

- refactor
- add tests
- go mod update

## v1.7.1

- MultiTrigger.Add returns Trigger instead Fire

## v1.7.0

- add ContextWithSig
- go mod update

## v1.6.0

- remove vendor
- go mod update

## v1.5.7

- go mod update

## v1.5.6

- go mod update

## v1.5.5

- go mod update

## v1.5.4

- go mod update

## v1.5.3

- go mod update
- replace pkg/errors

## v1.5.2

- go mod update

## v1.5.1

- return Func

## v1.5.0

- add backoff factor 

## v1.4.0

- add background runner

## v1.3.1

- use github.com/bborbe/errors for better error list display

## v1.3.0

- retry check if err is retryable
- update deps and add vulncheck

## v1.2.0

- use errors join

## v1.1.0

- use ginkgo v2
- improve use of counterfeiter

## v1.0.0

- Initial Version
