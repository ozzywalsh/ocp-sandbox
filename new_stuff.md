# OpenShift & OpenTelemetry Onboarding Task

This document tracks the progress and implementation details for the OpenShift local environment setup, observability configuration, and OpenTelemetry integration.

## 1. Environment Setup: OpenShift Local (CRC)
**Goal:** Install and verify a local OpenShift instance.

* **Documentation:** [CRC Installation Guide](https://crc.dev/docs/installing/)
* **Verification:** Run `crc status` to ensure the VM and cluster are operational.

```bash
# Target Status
CRC VM:          Running
OpenShift:       Running (v4.21.8)
...
```

---

## 2. Cluster Observability (Logging & Monitoring)
**Goal:** Enable User Workload Monitoring and set up the Loki/Vector logging stack.

### 📊 Monitoring
* **Guide:** [Configuring User Workload Monitoring](https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/monitoring/configuring-user-workload-monitoring)
* **Configuration:** Applied Cluster Monitoring ConfigMap: [`manifests/cluster/configmap-cluster-monitoring-config.yaml`](./manifests/cluster/configmap-cluster-monitoring-config.yaml)

### 🪵 Logging Stack (Loki & Vector)
To enable logging, the following components were deployed:
1.  **Loki Operator Subscription:** [`manifests/operator/subscription-loki.yaml`](./)
2.  **S3 Storage (SeaweedFS):** Initialized for Loki backend. [`manifests/storage/seaweedfs.yaml`](./)
3.  **LokiStack:** Configured with S3 secret and LokiStack CR. [`manifests/logging/lokistack.yaml`](./)



---

## 3. Application Metrics & OTel Integration
**Goal:** Transition from standard Prometheus scraping to an OTel-mediated pipeline.

### Phase A: Simple Instrumented Web App
* **App Source:** [Link to App Code](./todo-api/main.go)
* **Task:** Expose `/metrics` in Prometheus format.
* **Access:** Console credentials retrieved via `crc console --credentials`.

### Phase B: OpenTelemetry Operator & Collector
* **Operator:** Installed via [`manifests/operator/subscription-openshift-opentelemetry.yaml`](./)
* **Collector Logic:**
    * **Receiver:** `prometheus` (scrapes the app).
    * **Processors:** Added `resourcedetection` and `attributes` to inject K8s metadata and custom labels.
    * **Exporter:** `prometheusremotewrite` (back to OpenShift monitoring).

> [!IMPORTANT]
> **Files:** See the Collector CR definition here: [`manifests/workloads/otel-collector.yaml`](./)

---

## 4. Advanced Instrumentation: Tracing & Auto-Injection
**Goal:** Implement distributed tracing and move to "Zero-Code" instrumentation.

### Distributed Tracing
* **Operator:** Distributed Tracing (Tempo/Jaeger).
* **Implementation:** App updated to pass context headers across service boundaries.

### Auto-Instrumentation
* **Method:** Removed SDK libraries from the application code.
* **Mechanism:** Applied the `Instrumentation` CR and added the annotation to the Deployment:
    ```yaml
    instrumentation.opentelemetry.io/inject-sdk: "true"
    ```



---

## 5. Development & Contribution
**Goal:** Build and deploy custom versions of the OTel Operator and Collector.

### Local Build Process
1.  **Clone Repos:**
    * `git clone https://github.com/open-telemetry/opentelemetry-operator`
    * `git clone https://github.com/open-telemetry/opentelemetry-collector`
2.  **Build:** Followed `CONTRIBUTING.md` using `make build`.
3.  **Deployment:** Pushed local images to the internal OpenShift registry and patched the deployments.

| Component | Change Description | PR / Link |
| :--- | :--- | :--- |
| **Operator** | Added custom log line to controller start | [Link to Commit/PR]() |
| **Collector** | Modified version string for verification | [Link to Commit/PR]() |

---

### 📝 TODO List
- [ ] Link actual SeaweedFS manifest.
- [ ] Add screenshot of the Loki Logging UI in the OpenShift Console.
- [ ] Upload the OTLP-instrumented version of the application.
- [ ] Complete local build PR links.

---
