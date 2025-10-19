# Metrics and Monitoring

## Overview

This document describes the metrics, monitoring, and observability features of the k8s-pod-headless-service-operator.

## Current Logging

The operator uses [apex/log](https://github.com/apex/log) for structured logging:

- **Info logs**: Normal operation events (service creation, updates, deletions)
- **Error logs**: Failures and unexpected conditions

### Log Messages

The operator logs the following events:

1. **Pod Events**
   - `Watching pods` - Operator started successfully
   - `Setting up pod {name}` - Processing a new pod
   - `Updating pod {name}` - Processing a pod update
   - `Deleting service for pod {name}` - Processing a pod deletion

2. **Service Operations**
   - `{pod} doesn't have annotation set, skipping` - Pod doesn't match criteria
   - `{pod} doesn't have an IP yet skipping` - Pod has no IP assigned
   - `{pod} has a too long Pod name to create a service` - Name > 63 characters
   - `{pod} already has a service, updating it` - Service exists, will update
   - `{pod} has no service, creating it` - Creating new service
   - `{pod} has a new Pod IP, updating it` - Updating endpoint IP

3. **Errors**
   - `Error setting up service: {error}` - Failed to create service
   - `Error updating service: {error}` - Failed to update service
   - `Error deleting service: {error}` - Failed to delete service

## Code Coverage Metrics

The project tracks code coverage through:

- **Local testing**: `go test -coverprofile=coverage.out`
- **CI/CD**: Automated coverage reporting to Codecov
- **Current coverage**: 8.4% (baseline)

### Coverage Reports

Coverage reports are generated for:
- Unit tests in `cmd/k8s-pod-headless-service-operator/run_test.go`

View coverage locally:
```bash
go test -mod=mod ./cmd/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## CI/CD Metrics

GitHub Actions tracks:

- **Test pass/fail rates**: All PRs and commits
- **Build success rates**: Compilation success
- **Lint issues**: Code quality violations
- **Security scan results**: Potential vulnerabilities
- **YAML validation**: Manifest file correctness

## Future Metrics (Planned)

### Prometheus Metrics

Planned metrics to expose via `/metrics` endpoint:

```
# Services managed
k8s_headless_service_operator_services_total{namespace="default"} 10

# Service operations
k8s_headless_service_operator_operations_total{operation="create",status="success"} 25
k8s_headless_service_operator_operations_total{operation="update",status="success"} 50
k8s_headless_service_operator_operations_total{operation="delete",status="success"} 5
k8s_headless_service_operator_operations_total{operation="create",status="error"} 2

# Pod events processed
k8s_headless_service_operator_pod_events_total{event_type="add"} 30
k8s_headless_service_operator_pod_events_total{event_type="update"} 100
k8s_headless_service_operator_pod_events_total{event_type="delete"} 10

# Reconciliation time
k8s_headless_service_operator_reconcile_duration_seconds{quantile="0.5"} 0.05
k8s_headless_service_operator_reconcile_duration_seconds{quantile="0.9"} 0.15
k8s_headless_service_operator_reconcile_duration_seconds{quantile="0.99"} 0.5
```

### Health Checks

Planned health check endpoints:

- `/healthz` - Liveness probe
- `/readyz` - Readiness probe

## Monitoring Best Practices

1. **Watch logs** for error patterns
2. **Monitor service creation rates** - Unusual spikes may indicate issues
3. **Track annotation compliance** - Pods skipped due to missing annotations
4. **Alert on repeated failures** - Persistent errors need investigation

## Alerting (Future)

Recommended alerts:

```yaml
# High error rate
alert: HighServiceCreationErrorRate
expr: rate(k8s_headless_service_operator_operations_total{status="error"}[5m]) > 0.1
for: 10m
annotations:
  summary: High rate of service creation errors

# Operator down
alert: OperatorDown
expr: up{job="k8s-headless-service-operator"} == 0
for: 5m
annotations:
  summary: Operator is not running
```

## Dashboard (Future)

Grafana dashboard showing:

- Services managed over time
- Operation success/failure rates
- Event processing rates
- Reconciliation latency
- Error rate trends

## Tracing (Future)

Potential OpenTelemetry integration for:

- Request tracing through the operator
- Pod event processing flow
- Kubernetes API call latency
- Service creation/update/delete traces

## Current Observability Stack

For production deployments, integrate with:

- **Logs**: Collect via fluentd/fluent-bit to Elasticsearch or Loki
- **Metrics**: Export to Prometheus (when implemented)
- **Alerts**: Configure via Alertmanager (when metrics available)
- **Dashboards**: Create in Grafana (when metrics available)

## Debugging Metrics

For debugging issues:

1. Check pod annotations: `kubectl get pod <name> -o yaml | grep annotations -A5`
2. Check operator logs: `kubectl logs -n <namespace> <operator-pod>`
3. List services created: `kubectl get svc -l operator=headless-service-operator`
4. Check events: `kubectl get events --sort-by='.lastTimestamp'`

## Performance Metrics

Current performance characteristics:

- **Memory usage**: ~20-50MB baseline
- **CPU usage**: Minimal (<0.1 core typically)
- **API calls**: 1 watch + operations as needed
- **Reconciliation**: Event-driven (immediate response to pod changes)

## Contributing Metrics

When adding new metrics:

1. Use consistent naming (prefix: `k8s_headless_service_operator_`)
2. Include appropriate labels (namespace, operation, status)
3. Document in this file
4. Add to Grafana dashboard template
5. Consider adding corresponding alerts
