package controllers

import (
	"github.com/DoodleScheduling/k8skafka-controller/api/v1beta1"
	"github.com/DoodleScheduling/k8skafka-controller/kafka"
	"strconv"
)

// NOTE: kafka internal config values are intentionally not put in api.v1beta1 package, in order not to couple our API with kafka internal values.
// If kafka internals change between versions, we can create new translations without touching our API
func TranslateKafkaTopicV1Beta1(topic v1beta1.KafkaTopic) *kafka.Topic {
	c := kafka.Config{}
	// Boy would it be nice to have some generic configuration map[type]kafkaConfigString,
	// that can just go through each type (not concrete instance) in map, and get appropriate "translated" kafkaConfigString,
	// and get config value from concrete instance of that type
	//
	// something like (pseudo):
	// 		var v1betaKeys = map[type]string {
	//			v1beta1.CleanupPolicy:   "cleanup.policy",
	//			v1beta1.CompressionType: "compression.type",
	// 		}
	// and then we could just do:
	// 		for each config := v1betaKeys {
	//			c[config.value] = config.key.GetValue()
	//		}
	// Constraint being that kafka "names" like "cleanup.policy" shouldn't leak to api.v1beta1
	if topic.GetCleanupPolicy() != nil {
		c["cleanup.policy"] = string(*topic.GetCleanupPolicy())
	}
	if topic.GetCompressionType() != nil {
		c["compression.type"] = string(*topic.GetCompressionType())
	}
	if topic.GetDeleteRetentionMs() != nil {
		c["delete.retention.ms"] = strconv.FormatInt(*topic.GetDeleteRetentionMs(), 10)
	}
	if topic.GetFileDeleteDelayMs() != nil {
		c["file.delete.delay.ms"] = strconv.FormatInt(*topic.GetFileDeleteDelayMs(), 10)
	}
	if topic.GetFlushMessages() != nil {
		c["flush.messages"] = strconv.FormatInt(*topic.GetFlushMessages(), 10)
	}
	if topic.GetFlushMs() != nil {
		c["flush.ms"] = strconv.FormatInt(*topic.GetFlushMs(), 10)
	}
	if topic.GetFollowerReplicationThrottledReplicas() != nil {
		c["follower.replication.throttled.replicas"] = *topic.GetFollowerReplicationThrottledReplicas()
	}
	if topic.GetIndexIntervalBytes() != nil {
		c["index.interval.bytes"] = strconv.FormatInt(*topic.GetIndexIntervalBytes(), 10)
	}
	if topic.GetLeaderReplicationThrottledReplicas() != nil {
		c["leader.replication.throttled.replicas"] = *topic.GetLeaderReplicationThrottledReplicas()
	}
	if topic.GetMaxMessageBytes() != nil {
		c["max.message.bytes"] = strconv.FormatInt(*topic.GetMaxMessageBytes(), 10)
	}
	if topic.GetMessageDownconversionEnable() != nil {
		c["message.downconversion.enable"] = strconv.FormatBool(*topic.GetMessageDownconversionEnable())
	}
	if topic.GetMessageFormatVersion() != nil {
		c["message.format.version"] = *topic.GetMessageFormatVersion()
	}
	if topic.GetMessageTimestampDifferenceMaxMs() != nil {
		c["message.timestamp.difference.max.ms"] = strconv.FormatInt(*topic.GetMessageTimestampDifferenceMaxMs(), 10)
	}
	if topic.GetMessageTimestampType() != nil {
		c["message.timestamp.type"] = string(*topic.GetMessageTimestampType())
	}
	if topic.GetMinCleanableDirtyRatio() != nil {
		c["min.cleanable.dirty.ratio"] = strconv.FormatInt(*topic.GetMinCleanableDirtyRatio(), 10)
	}
	if topic.GetMinCompactionLagMs() != nil {
		c["min.compaction.lag.ms"] = strconv.FormatInt(*topic.GetMinCompactionLagMs(), 10)
	}
	if topic.GetMinInsyncReplicas() != nil {
		c["min.insync.replicas"] = strconv.FormatInt(*topic.GetMinInsyncReplicas(), 10)
	}
	if topic.GetPreallocate() != nil {
		c["preallocate"] = strconv.FormatBool(*topic.GetPreallocate())
	}
	if topic.GetRetentionBytes() != nil {
		c["retention.bytes"] = strconv.FormatInt(*topic.GetRetentionBytes(), 10)
	}
	if topic.GetRetentionMs() != nil {
		c["retention.ms"] = strconv.FormatInt(*topic.GetRetentionMs(), 10)
	}
	if topic.GetSegmentBytes() != nil {
		c["segment.bytes"] = strconv.FormatInt(*topic.GetSegmentBytes(), 10)
	}
	if topic.GetSegmentIndexBytes() != nil {
		c["segment.index.bytes"] = strconv.FormatInt(*topic.GetSegmentIndexBytes(), 10)
	}
	if topic.GetSegmentJitterMs() != nil {
		c["segment.jitter.ms"] = strconv.FormatInt(*topic.GetSegmentJitterMs(), 10)
	}
	if topic.GetSegmentMs() != nil {
		c["segment.ms"] = strconv.FormatInt(*topic.GetSegmentMs(), 10)
	}
	if topic.GetUncleanLeaderElectionEnable() != nil {
		c["unclean.leader.election.enable"] = strconv.FormatBool(*topic.GetUncleanLeaderElectionEnable())
	}

	kt := kafka.Topic{
		Name:              topic.GetTopicName(),
		Partitions:        topic.GetPartitions(),
		ReplicationFactor: topic.GetReplicationFactor(),
		Config:            c,
	}

	return &kt
}
