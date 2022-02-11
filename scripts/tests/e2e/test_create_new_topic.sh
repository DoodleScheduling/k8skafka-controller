#!/bin/bash
set -e

echo "test creating new topic"

# create topic
kubectl -n k8skafka-system apply -f ./config/testdata/kafkatopics/create-new
kubectl -n k8skafka-system wait kafkatopic/test-create-new --for=condition=Ready --timeout=5m

# get topic from Kafka
res=$(kubectl -n kafka exec -i kafka-client -- kafka-topics --describe --topic test-create-new --bootstrap-server kafka-cp-kafka:9092)

# assert topic exists in Kafka
[[ -z "$res" ]] && (there is no topic test-create-new && exit 1) || echo "topic test-create-new successfully created"