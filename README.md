# Kubernetes Pod Headless Service Operator

[![CI](https://github.com/src-d/k8s-pod-headless-service-operator/workflows/CI/badge.svg)](https://github.com/src-d/k8s-pod-headless-service-operator/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/src-d/k8s-pod-headless-service-operator)](https://goreportcard.com/report/github.com/src-d/k8s-pod-headless-service-operator)
[![codecov](https://codecov.io/gh/src-d/k8s-pod-headless-service-operator/branch/master/graph/badge.svg)](https://codecov.io/gh/src-d/k8s-pod-headless-service-operator)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This is a Kubernetes operator that watches Pods with the annotation `srcd.host/create-headless-service: "true"`, if this annotation is found the operator will create a Headless Service with the pod's name and an Endpoint pointing to the Pod IP. This allows for the Pod's hostname to be resolvable by DNS, this is a requirement needed by certain applications.

## Features

- Automatically creates headless services for annotated pods
- Manages endpoints to ensure DNS resolution works correctly
- Updates endpoints when pod IPs change
- Cleans up services when pods are deleted
- Namespace filtering support
- Configurable annotation key

## Limitations
This will only work if your Pod's name is maximum 63 characters as this is the maximum length for a service name.

# Installation

This tool is made to run in cluster as a Deployment. For testing purposes it can also run locally with a connection to a Kubernetes cluster.

## Kubernetes manifests
This repository provides example manifests file you can use to deploy this. These contain a service account and RBAC configuration for the tool to be able to read Pods and read/write Services. As well as a Deployment to deploy the operator in a cluster.
```bash
~ $ cd manifests
~ $ kubectl apply -f rbac.yaml
~ $ kubectl apply -f daemonset.yaml
```

## Helm
We also provide a Helm chart in our [Charts repository](https://github.com/src-d/charts). 
```bash
~ $ helm repo add srcd-infra https://src-d.github.io/charts/infra/
~ $ helm install k8s-pod-headless-service-operator --set image.tag=v0.1.1
```

# Configuration

* envvar: `NAMESPACE` flag: `--namespace` The namespace to watch, by default it watches all namespaces
* envvar: `POD_ANNOTATION` flag: `--pod-annotation` Pod annotation that needs to be set to `true` to be picked up by the operator. Default: `srcd.host/create-headless-service`
* envvar: `KUBERNETES_CONTEXT` flag: `--context` If this is set it will not attempt to load the in-cluster service account but loads the context value out of `$HOME/.kube/config`

# Development

For development instructions, see [DEVELOPMENT.md](DEVELOPMENT.md).

For testing guide, see [TESTING.md](TESTING.md).

For metrics and monitoring, see [METRICS.md](METRICS.md).

## Quick Start for Developers

```bash
# Clone the repository
git clone https://github.com/src-d/k8s-pod-headless-service-operator.git
cd k8s-pod-headless-service-operator

# Download dependencies
go mod download

# Run tests
go test -mod=mod ./cmd/... -v

# Build
go build -mod=mod -v ./cmd/k8s-pod-headless-service-operator
```

# Contribute

[Contributions](https://github.com/src-d/k8s-pod-headless-service-operator/issues) are more than welcome, if you are interested please take a look to
our [Contributing Guidelines](CONTRIBUTING.md).

# Code of Conduct

All activities under source{d} projects are governed by the [source{d} code of conduct](.github/CODE_OF_CONDUCT.md).

# License
Apache License Version 2.0, see [LICENSE](LICENSE).
