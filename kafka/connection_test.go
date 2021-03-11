package kafka

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssignBrokers_brokersAreUnique(t *testing.T) {
	brokers := make(map[int64]int64)
	brokers[1] = 16
	brokers[2] = 16
	brokers[3] = 16
	for i := 0; i < 1000; i++ {
		conn := NewTCPConnection("kafka:9092")
		topic := Topic{
			Name:              "testTopic",
			Partitions:        16,
			ReplicationFactor: 3,
			Brokers:           brokers,
		}
		var brokerIDs []int32
		brokerIDs, brokers = conn.assignBrokersToPartition(topic, brokers)

		u := unique(brokerIDs)
		assert.Equal(t, 3, len(u))
	}
}

func unique(intSlice []int32) []int32 {
	keys := make(map[int32]bool)
	var list []int32
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
