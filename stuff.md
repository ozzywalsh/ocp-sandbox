# Explore OpenShift

* Install OpenShift local on your company provided laptop
You can install CRC using the guide here: https://crc.dev/docs/installing/
You can start CRC using optimal resources for the observability stack using the script in  [`scripts/start-crc`](./)
You can verify the installation by running `crc status`
```
ozwalsh@ozwalsh-thinkpadt14gen4:~$ crc status
CRC VM:          Running
OpenShift:       Running (v4.21.8)
RAM Usage:       19.07GB of 25.19GB
Disk Usage:      41.65GB of 267.8GB (Inside the CRC VM)
Cache Usage:     94.2GB
Cache Directory: /home/ozwalsh/.crc/cache
```

* Ensure that OpenShift local has monitoring enabled (in-cluster and user workload monitoring) and that it is configured to use logging (Vector & Loki)
OpenShift includes core platform monitoring out of the box[^1].

To enable the user workload monitoring follow this guide: https://docs.r>edhat.com/en/documentation/openshift_container_platform/4.14/html/monitoring/configuring-user-workload-monitoring
For an example check `./manifests/cluster/configmap-cluster-monitoring-config.yaml`
TODO: Add link to logging docs
To do the installation on your CRC cluster we first need to apply the subscription to install the loki-operator. TODO: Link subscription
Next we need to prepare an s3 compatible backend for loki. This repo utilizes seaweedfs for this purpose.
1. Seaweedfs TODO: Link to file
2. Create a secret containing the credentials for s3 TODO: link
3. Finally create a LokiStack TODO: add link to file
Your `openshift-logging` namespace should look like this:
```
ozwalsh@ozwalsh-thinkpadt14gen4:~$ kubectl -n openshift-logging get pods
NAME                                          READY   STATUS    RESTARTS       AGE
cluster-logging-operator-77bfb75c76-vghlt     1/1     Running   2              2d
instance-xdxtg                                1/1     Running   2              2d
logging-loki-compactor-0                      1/1     Running   11 (21h ago)   2d
logging-loki-distributor-748b788dcd-m7zv8     1/1     Running   2              2d
logging-loki-gateway-6696d6f469-kjshc         2/2     Running   4              2d
logging-loki-gateway-6696d6f469-rgbxm         2/2     Running   4              2d
logging-loki-index-gateway-0                  1/1     Running   2              2d
logging-loki-ingester-0                       1/1     Running   2              2d
logging-loki-querier-76997b458b-j4c6d         1/1     Running   2              2d
logging-loki-query-frontend-ff69bc95f-rpqdw   1/1     Running   2              2d
```

* Create a simple web application that will expose custom metrics in the Prometheus format and output logs.
TODO add link to prometheus instrumented app

* Ensure you can access logs and custom metrics in the OpenShift console for your application. For logs, ensure you are using the Loki based logging console page.
First install the cluster observability-operator by adding the relevant subscription, operatorgroup etc. `./manifests/operator/subscription-cluster-observability.yaml`
Install the logging UI plugin by applying the following resource to your cluster `./manifests/operator/uiplugin-logging.yaml`
To get the admin credentials for the web console run:
```
crc console --credentials
```
Then use the kubadmin user and the output password to login here: https://console-openshift-console.apps-crc.testing/

TODO: Add screenshot
* Install OpenTelemetry
To install the opentelemetry-operator apply the opentelemetry subscription here `./manifests/operator/subscription-openshift-opentelemetry-product.yaml`
This should result in:
```
ozwalsh@ozwalsh-thinkpadt14gen4:~/dev/otel-sndbox$ kubectl -n openshift-opentelemetry-operator get pods
NAME                                                         READY   STATUS    RESTARTS      AGE
opentelemetry-operator-controller-manager-86cbd7f9dc-hfpdk   1/1     Running   6 (21h ago)   2d
```

* Configure your OpenTelemetry collector so that it will read the metrics endpoint instead of the OpenShift Prometheus (the metrics endpoint should still be in the prometheus format, but scraped by the OTEL collector and not by the OpenShift Prometheus). Add an extra label to your metrics using the Collector and make sure that Collector is adding the right metadata for an application running in OpenShift/Kubernetes (you may need a processor to help with this). Ensure that you can see your metrics (with extra labels) in the OpenShift console.
To fulfill this step we must first create a `OpenTelemetryCollector` CR configured with the `prometheus` receiver. 
We must also add a processor add the appropriate labels to the metrics. See file `./manifests/workloads/opentelemetrycollector-otel`. (TODO: link old commit with prometheus receiver)

* Update your application so that it uses OTLP metrics and logs. This will require updating your application to use the OpenTelemetry instrumentation libraries instead of the Prometheus one. Both the logs and metrics should be going through the OpenTelemetry collector. Ensure that you can see these metrics and logs in the OpenShift Console.

* Install Distributed Tracing in the OpenShift console and update your simple application so that a request goes through multiple different services.

Instrument your application for tracing (autoinstrumentation is fine, but please do so at build time).
Ensure that you can now access your OTEL metrics, logs and traces in the OpenShift console.
Remove the instrumentation libraries for your application. This time use the OpenTelemetry Operator to inject auto-intrumentation into your application. Ensure you can access the logs, metrics and traces in the OpenShift Console.


* Checkout the OpenTelemetry Operator and Collector code bases. Ensure that you can build these components locally on your machine.
Clone this repo (otel repo). Follow the contributing.md. Link pull requests

Deploy these locally built products into your OpenShift Local cluster.
Make a small change to the Operator and Collector to ensure that you are deploying this new version.



[^1]: https://docs.redhat.com/en/documentation/monitoring_stack_for_red_hat_openshift/4.21/html/about_monitoring/about-ocp-monitoring