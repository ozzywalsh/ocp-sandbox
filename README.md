# OpenTelemetry Sandbox

OpenShift (CRC) cluster configuration for deploying the OpenTelemetry operator, a collector instance, and monitoring infrastructure.

## Structure

```
manifests/
  cluster/          # Cluster-scoped resources: namespaces, RBAC, API server config, alertmanager
  operator/         # OLM resources for installing the OpenTelemetry operator
  workloads/        # Application resources: collectors, monitoring config, mailpit
```

## Apply Order

```bash
# 1. Cluster-scoped resources (namespaces, RBAC, cluster config, alertmanager)
kubectl apply -f manifests/cluster/

# 2. Install the OpenTelemetry operator
kubectl apply -f manifests/operator/

# 3. Wait for the operator to be ready
kubectl wait --for=condition=Available deployment/opentelemetry-operator-controller-manager \
  -n openshift-opentelemetry-operator --timeout=120s

# 4. Deploy workloads
kubectl apply -f manifests/workloads/
```

## MailPit

MailPit runs in the `sandbox` namespace as a local SMTP server for Alertmanager notifications. Alertmanager is configured to send all alerts (except Watchdog) to MailPit.

Access the web UI via the OpenShift route:

```bash
oc get route mailpit -n sandbox
# opens at http://mailpit-sandbox.apps-crc.testing
```
