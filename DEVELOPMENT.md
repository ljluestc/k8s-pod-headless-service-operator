# Development Guide

## Prerequisites

- Go 1.23 or later
- Docker (for building container images)
- kubectl (for testing in Kubernetes clusters)
- Access to a Kubernetes cluster (optional, for integration testing)

## Setting Up Development Environment

1. **Clone the repository**

```bash
git clone https://github.com/src-d/k8s-pod-headless-service-operator.git
cd k8s-pod-headless-service-operator
```

2. **Install dependencies**

```bash
go mod download
```

3. **Build the project**

```bash
go build -mod=mod -v ./cmd/k8s-pod-headless-service-operator
```

## Development Workflow

### Running Tests

See [TESTING.md](TESTING.md) for detailed testing instructions.

```bash
# Run all tests
go test -mod=mod ./cmd/... -v

# Run tests with coverage
go test -mod=mod ./cmd/... -v -coverprofile=coverage.out

# Run tests with race detection
go test -mod=mod ./cmd/... -v -race
```

### Code Linting

```bash
# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest

# Run linter
golangci-lint run --timeout=5m
```

### Building

```bash
# Build binary
go build -mod=mod -o bin/k8s-pod-headless-service-operator ./cmd/k8s-pod-headless-service-operator

# Build Docker image
docker build -t k8s-pod-headless-service-operator:dev .
```

## Project Structure

```
.
├── cmd/
│   └── k8s-pod-headless-service-operator/
│       ├── main.go           # Application entry point
│       ├── run.go            # Main operator logic
│       └── run_test.go       # Unit tests
├── manifests/
│   ├── rbac.yaml             # RBAC configuration
│   └── deployment.yaml       # Kubernetes deployment
├── .github/
│   └── workflows/
│       └── ci.yml            # GitHub Actions CI/CD
├── vendor/                   # Vendored dependencies
├── .golangci.yml             # Linter configuration
├── Dockerfile                # Container image definition
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── Makefile                  # Build automation
├── README.md                 # Project overview
├── TESTING.md                # Testing guide
├── DEVELOPMENT.md            # This file
└── CONTRIBUTING.md           # Contribution guidelines
```

## Running Locally

### Against a Local Cluster

```bash
# Set the Kubernetes context
export KUBERNETES_CONTEXT=minikube  # or your context name

# Run the operator
./bin/k8s-pod-headless-service-operator run --context=$KUBERNETES_CONTEXT
```

### Configuration Options

The operator supports the following configuration options:

- `--namespace` / `NAMESPACE`: Namespace to watch (default: all namespaces)
- `--pod-annotation` / `POD_ANNOTATION`: Annotation to look for (default: `srcd.host/create-headless-service`)
- `--context` / `KUBERNETES_CONTEXT`: Kubernetes context to use (for local development)

## Debugging

### Enable Verbose Logging

The operator uses the apex/log library. You can control logging verbosity through environment variables.

### Common Issues

1. **Tests fail with protobuf conflicts**: The project has specific version constraints for protobuf libraries. See `go.mod` replace directives.

2. **Vendor directory inconsistency**: Run `go mod vendor` or use `-mod=mod` flag when running go commands.

3. **RBAC errors in cluster**: Ensure the ServiceAccount has proper permissions defined in `manifests/rbac.yaml`.

## Code Quality Standards

- All code must pass `golangci-lint` checks
- New features must include tests
- Code coverage should not decrease
- Follow Go best practices and idioms
- Use meaningful variable and function names
- Add comments for exported functions and types

## Making Changes

1. Create a feature branch from `master`
2. Make your changes
3. Add or update tests
4. Run tests and linter locally
5. Commit with clear, descriptive messages
6. Push and create a pull request
7. Wait for CI checks to pass
8. Address review feedback

## Release Process

Releases are automated through GitHub Actions when tags are pushed:

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

This will:
- Run all CI checks
- Build Docker images
- Push images to Docker Hub

## Getting Help

- Check existing [Issues](https://github.com/src-d/k8s-pod-headless-service-operator/issues)
- Review [Contributing Guidelines](CONTRIBUTING.md)
- Join community discussions

## Resources

- [Kubernetes Client-Go Documentation](https://github.com/kubernetes/client-go)
- [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Go Testing](https://golang.org/pkg/testing/)
