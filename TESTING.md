# Testing Guide

## Overview

This document describes the testing infrastructure for the k8s-pod-headless-service-operator project.

## Test Suite

The project includes comprehensive unit tests that verify the core functionality of the operator:

- **Service Creation Logic**: Tests for creating headless services for annotated pods
- **Service Update Logic**: Tests for updating service endpoints when pod IPs change
- **Service Deletion Logic**: Tests for cleaning up services when pods are deleted
- **Annotation Handling**: Tests for proper annotation detection and validation
- **Name Length Validation**: Tests for handling pods with names exceeding 63 characters

## Running Tests

### Run all tests

```bash
go test -mod=mod ./cmd/... -v
```

### Run tests with coverage

```bash
go test -mod=mod ./cmd/... -v -coverprofile=coverage.out -covermode=atomic
```

### View coverage report

```bash
go tool cover -html=coverage.out
```

### Run tests with race detection

```bash
go test -mod=mod ./cmd/... -v -race
```

## Test Structure

Tests are located in `cmd/k8s-pod-headless-service-operator/run_test.go` and use table-driven testing patterns for comprehensive coverage.

### Key Test Functions

1. **TestGetClientSet**: Tests Kubernetes client initialization
2. **TestHasExistingServiceLogic**: Tests service existence checking
3. **TestServiceCreationLogic**: Tests the logic for determining when to create services
4. **TestServiceCreationWithFakeClient**: Tests actual service and endpoint creation
5. **TestEndpointUpdate**: Tests endpoint IP address updates
6. **TestServiceDeletion**: Tests service cleanup

## Continuous Integration

The project uses GitHub Actions for continuous integration with the following checks:

- Unit tests with coverage reporting
- Code linting (golangci-lint)
- Security scanning (gosec)
- YAML linting
- Build verification

See `.github/workflows/ci.yml` for the full CI configuration.

## Code Coverage

Current coverage target: 8.4% (baseline)

Coverage reports are automatically uploaded to Codecov on each CI run.

## Testing Best Practices

1. **Use Fake Kubernetes Clients**: All tests use `k8s.io/client-go/kubernetes/fake` to avoid requiring a real cluster
2. **Table-Driven Tests**: Most tests use table-driven patterns for maintainability
3. **Clear Test Names**: Test names describe what is being tested and expected behavior
4. **Isolated Tests**: Each test is independent and doesn't rely on shared state

## Adding New Tests

When adding new functionality:

1. Add corresponding test cases in `run_test.go`
2. Use the existing test patterns (table-driven tests with fake clients)
3. Run tests locally before committing
4. Ensure CI passes before merging

## Dependencies

Test dependencies are managed through `go.mod` and include:

- `k8s.io/client-go/kubernetes/fake`: Fake Kubernetes clients for testing
- `k8s.io/api/core/v1`: Kubernetes core API types
- `k8s.io/apimachinery/pkg/apis/meta/v1`: Kubernetes metadata types
