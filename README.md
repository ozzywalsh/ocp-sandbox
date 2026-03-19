# OpenTelemetry Sandbox

OpenShift cluster configuration for deploying the OpenTelemetry operator and a collector instance.

## Structure

```
manifests/
  cluster/          # Cluster-scoped resources: namespaces, RBAC, API server config
  operator/         # OLM resources for installing the OpenTelemetry operator
  workloads/        # Application resources: collectors, deployments, services
```

## Apply Order

Resources must be applied in stages since the operator needs to be running before its CRDs can be used.

```bash
# 1. Cluster-scoped resources (namespaces, RBAC, cluster config)
kubectl apply -f manifests/cluster/

# 2. Install the OpenTelemetry operator
kubectl apply -f manifests/operator/

# 3. Wait for the operator to be ready
kubectl wait --for=condition=Available deployment/opentelemetry-operator-controller-manager \
  -n openshift-opentelemetry-operator --timeout=120s

# 4. Deploy workloads
kubectl apply -f manifests/workloads/
```
