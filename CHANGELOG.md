# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive unit test suite with 6 test functions covering core functionality
- GitHub Actions CI/CD workflow with multiple jobs:
  - Unit tests with coverage reporting
  - Code linting with golangci-lint
  - Security scanning with gosec
  - YAML validation
  - Build verification
- Go modules support (migrated from dep)
- Test coverage reporting to Codecov
- golangci-lint configuration for code quality
- Comprehensive documentation:
  - TESTING.md - Testing guide
  - DEVELOPMENT.md - Development workflow guide
  - METRICS.md - Metrics and monitoring guide
  - CHANGELOG.md - This file

### Changed
- Updated to Go 1.23
- Migrated from dep to Go modules for dependency management
- Fixed protobuf dependency conflicts
- Updated TravisCI configuration to use modern Go versions

### Fixed
- Resolved inconsistencies between vendor directory and go.mod
- Fixed protobuf version conflicts with proper replace directives
- Ensured all tests pass with fake Kubernetes clients

## [0.1.1] - Previous Release

### Added
- Initial operator implementation
- Pod watching functionality
- Headless service creation for annotated pods
- Endpoint management for pod IPs
- RBAC manifests
- Deployment manifests
- Helm chart support

### Features
- Watches pods with annotation `srcd.host/create-headless-service: "true"`
- Creates headless services (ClusterIP: None) for matching pods
- Creates endpoints pointing to pod IPs
- Updates endpoints when pod IPs change
- Deletes services when pods are removed
- Namespace filtering support
- Configurable annotation key

## [0.1.0] - Initial Release

### Added
- Basic operator structure
- Kubernetes client integration
- Pod informer setup
- Service and endpoint creation logic
