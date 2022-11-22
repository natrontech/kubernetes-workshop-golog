# Monitoring

Monitoring is a key part of any production system. It allows you to keep track of the health of your system and to identify problems before they become critical. Kubernetes provides a number of tools that can help you monitor your system.

You can also checkout the [Client Tools](../setup/tools/index.md) page for a list of local tools, which can help you monitor your cluster.

## Kubernetes Dashboard

- [Kubernetes Dashboard](https://dashboard.k8s.golog.ch)

The Kubernetes Dashboard is a web-based Kubernetes user interface. You can use the Kubernetes Dashboard to deploy containerized applications to a Kubernetes cluster, troubleshoot your containerized application, and manage the cluster resources. You can use the Kubernetes Dashboard to get an overview of applications running on your cluster, as well as for creating or modifying individual Kubernetes resources (such as Deployments, Jobs, DaemonSets, etc). For example, you can scale a Deployment, initiate a rolling update, restart a pod or deploy new applications using a deploy wizard.

The Dashboard also provides information on the state of Kubernetes resources in your cluster and on any errors that may have occurred.

### Metrics Server

The Kubernetes Dashboard UI (in particular, the graphs on the *Cluster* page) displays resource usage data (CPU and memory) of your cluster's nodes. The resource usage data is retrieved by the Kubernetes Dashboard from the Metrics Server. The Metrics Server is not deployed by default in Kubernetes. See the [Metrics Server documentation](https://github.com/kubernetes-sigs/metrics-server) to learn how to deploy the Metrics Server.

You have to redeploy the Kubernetes Dashboard with additional permissions to allow it to access the Metrics Server.

## Kube Prometheus Stack

The [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack) is a collection of community-curated Helm charts that deploy the core components of the Prometheus monitoring system. The kube-prometheus-stack Helm chart deploys the following components:

- [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator)
- [Prometheus](https://prometheus.io/)
- [Alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager/)
- [Prometheus Node Exporter](https://github.com/prometheus/node_exporter)
- [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)
- [Grafana](https://grafana.com/)

### Installation

To install the kube-prometheus-stack Helm chart, run the following command:

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
kubectl create namespace monitoring
helm install monitoring prometheus-community/kube-prometheus-stack -n monitoring
kubectl get pods -n monitoring -w
```

Have a look at the `values.yaml` file to see all the configuration options.
The default configuration should be fine for most use cases, but you can add ingress rules to access the services from outside the cluster.

### Accessing the Grafana Dashboard

To access the Grafana dashboard, run the following command:

```bash
kubectl port-forward svc/monitoring-grafana 3000:80 -n monitoring
```

Then, open the following URL in your browser: [http://localhost:3000](http://localhost:3000)

Search for the `monitoring-grafana` secret and base64 decode the `admin-user` and `admin-password` values.

```bash
kubectl get secret monitoring-grafana -n monitoring -o jsonpath="{.data.admin-user}" | base64 --decode ; echo
kubectl get secret monitoring-grafana -n monitoring -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

!!! note

    The default username is `admin` and the default password is `prom-operator`.

### Accessing the Prometheus Dashboard

To access the Prometheus dashboard, run the following command:

```bash
kubectl port-forward svc/monitoring-kube-prometheus-prometheus 9090:9090 -n monitoring
```

Then, open the following URL in your browser: [http://localhost:9090](http://localhost:9090)

## Loki

- [Loki](https://grafana.com/oss/loki/)

Loki is a horizontally-scalable, highly-available, multi-tenant log aggregation system inspired by Prometheus. It is designed to be very cost effective and easy to operate. It does not index the contents of the logs, but rather a set of labels for each log stream.
