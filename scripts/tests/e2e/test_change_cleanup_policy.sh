#!/bin/bash
set -e

echo "test change cleanup policy"

# change cleanup policy
kubectl -n k8skafka-system apply -f ./config/testdata/kafkatopics/change-cleanup-policy
sleep 2

# expected result
expectedPolicy="cleanup.policy=compact,delete"

# error variable holder, initially assume test will fail
e=true
# get expected result attempts iterator
i=1
# number of tries to get expected result
numberOfAttempts=5
# wait in seconds between attempts
timeoutForAttempt=3

while [ $i -le $numberOfAttempts ]; do
  # get result of changing cleanup policy in Kafka
  res=$(kubectl -n k8skafka-system exec -i sts/kafka -- kafka-topics.sh --describe --topic test-create-new --bootstrap-server kafka:9092)
  pol=$(echo "$res" | grep "$expectedPolicy" | wc -l | xargs)

  # assert policy is as expected in Kafka
  if [[ $pol == 1 ]]; then
    echo "topic has expected policy $expectedPolicy"
    e=false
    break
  else
    echo "expected to have $expectedPolicy. Full topic: $res"
    (( "i+=1" ))
    sleep $timeoutForAttempt
  fi
done

# fail if cleanup policy in Kafka is not as expected after timeout
if [[ $e == true ]]; then
  echo "failed to get required cleanup policy in time"
  exit 1
fi
