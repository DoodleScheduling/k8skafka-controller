#!/bin/bash
set -e

echo "test adding new partitions"

# add partitions
kubectl -n k8skafka-system apply -f ./config/testdata/kafkatopics/add-partitions
sleep 2

# expected result
expectedPartitions=20

# error variable holder, initially assume test will fail
e=true
# get expected result attempts iterator
i=1
# number of tries to get expected result
numberOfAttempts=5
# wait in seconds between attempts
timeoutForAttempt=3

while [ $i -le $numberOfAttempts ]; do
  # get result of adding partitions in Kafka
  res=$(kubectl -n k8skafka-system exec -i sts/kafka -- kafka-topics.sh --describe --topic test-create-new --bootstrap-server kafka:9092)
  res=$(echo "$res" | tail -n+2 | wc -l | xargs)

  # assert number of partitions in Kafka is as expected
  if [[ $res == "$expectedPartitions" ]]; then
    echo "topic has expected number of partitions: $expectedPartitions"
    e=false
    break
  else
    echo "expected $expectedPartitions partitions, got $res"
    (( "i+=1" ))
    sleep $timeoutForAttempt
  fi
done

# fail if number of partitions in Kafka is not as expected after timeout
if [[ $e == true ]]; then
  echo "failed to get required number of partitions in time"
  exit 1
fi
