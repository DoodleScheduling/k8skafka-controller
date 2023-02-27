# k8skafka-controller

[![release](https://img.shields.io/github/release/DoodleScheduling/k8skafka-controller/all.svg)](https://github.com/DoodleScheduling/k8skafka-controller/releases)
[![release](https://github.com/doodlescheduling/k8skafka-controller/actions/workflows/release.yaml/badge.svg)](https://github.com/doodlescheduling/k8skafka-controller/actions/workflows/release.yaml)
[![report](https://goreportcard.com/badge/github.com/DoodleScheduling/k8skafka-controller)](https://goreportcard.com/report/github.com/DoodleScheduling/k8skafka-controller)
[![Coverage Status](https://coveralls.io/repos/github/DoodleScheduling/k8skafka-controller/badge.svg?branch=master)](https://coveralls.io/github/DoodleScheduling/k8skafka-controller?branch=master)
[![license](https://img.shields.io/github/license/DoodleScheduling/k8skafka-controller.svg)](https://github.com/DoodleScheduling/k8skafka-controller/blob/master/LICENSE)

Kubernetes controller that can manage Kafka Topics.

Controller assumes you already have a working Kafka installation. How Kafka brokers are actually managed is outside of the scope of this project for the moment.

## Features
### Supported
- creating topics
- creating additional partitions for existing topics

### Not currently supported
- deleting topics, garbage collection
- deleting partitions
- changing replication factor for existing topics and partitions

### Known issues
- partition assigning can in some cases lead to leader being skewed

## Example KafkaTopic

A `KafkaTopic` represents one Kafka Topic.

```yaml
apiVersion: kafka.infra.doodle.com/v1beta1
kind: KafkaTopic
metadata:
  name: test-topic
spec:
  address: "kafka:9092"
  name: "test-topic"
  partitions: 16
  replicationFactor: 1
  config:
    cleanupPolicy: compact
    deleteRetentionsMs: 31556952000
```

## Helm chart

Please see [chart/k8skafka-controller](https://github.com/DoodleScheduling/k8skafka-controller/tree/master/chart/k8skafka-controller) for the helm chart docs.

## Configure the controller

You may change base settings for the controller using env variables (or alternatively command line arguments).
Available env variables:

| Name  | Description | Default |
|-------|-------------| --------|
| `METRICS_ADDR` | The address of the metric endpoint binds to. | `:9556` |
| `PROBE_ADDR` | The address of the probe endpoints binds to. | `:9557` |
| `ENABLE_LEADER_ELECTION` | Enable leader election for controller manager. | `false` |
| `LEADER_ELECTION_NAMESPACE` | Change the leader election namespace. This is by default the same where the controller is deployed. | `` |
| `NAMESPACES` | The controller listens by default for all namespaces. This may be limited to a comma delimited list of dedicated namespaces. | `` |
| `CONCURRENT` | The number of concurrent reconcile workers.  | `4` |
