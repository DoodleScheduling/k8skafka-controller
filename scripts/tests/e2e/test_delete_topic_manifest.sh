#!/bin/bash
set -e

echo "test delete topic"

# delete topic manifest
kubectl -n k8skafka-system delete kafkatopic/test-create-new
sleep 2

# expected result
expectedResult="NotFound"

# error variable holder, initially assume test will fail
e=true
# get expected result attempts iterator
i=1
# number of tries to get expected result
numberOfAttempts=5
# wait in seconds between attempts
timeoutForAttempt=3

while [ $i -le $numberOfAttempts ]; do
  # get result of deleting topic manifest
  res=$(kubectl -n k8skafka-system get kafkatopic/test-create-new || echo "NotFound")

  # assert topic manifest no longer exists
  if [[ "$res" == "$expectedResult" ]]; then
    e=false
    break
  else
    echo "KafkaTopic test-create-new not yet deleted. Res: $res"
    (( "i+=1" ))
    sleep $timeoutForAttempt
  fi
done

# fail if topic manifest still exists after timeout
if [[ $e == true ]]; then
  echo "failed to delete KafkaTopic in time"
exit 1
fi

# assert that topic still exists in Kafka, since we don't support garbage collection for now
res=$(kubectl -n kafka exec -i kafka-client -- kafka-topics --describe --topic test-create-new --bootstrap-server kafka-cp-kafka:9092)
[[ -z "$res" ]] && (there is no topic test-create-new && exit 1) || echo "topic test-create-new still exists, as expected"