# k8skafka controller helm chart

Installs the [k8skafka-controller](https://github.com/DoodleScheduling/k8skafka-controller).

## Installing the Chart

To install the chart with the release name `k8skafka-controller`:

```console
helm upgrade --install k8skafka-controller chart/k8skafka-controller
```

This command deploys the k8skafka-controller with the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Using the Chart

The chart comes with a ServiceMonitor for use with the [Prometheus Operator](https://github.com/helm/charts/tree/master/stable/prometheus-operator).
If you're not using the Prometheus Operator, you can disable the ServiceMonitor by setting `serviceMonitor.enabled` to `false` and instead
populate the `podAnnotations` as below:

```yaml
podAnnotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "metrics"
  prometheus.io/path: "/metrics"
```

## Configuration

See Customizing the Chart Before Installing. To see all configurable options with detailed comments, visit the chart's values.yaml, or run the configuration command:

```sh
$ helm show values k8skafka-controller/k8skafka-controller
```
