package controllers

import (
	"github.com/DoodleScheduling/k8skafka-controller/api/v1beta1"
	"github.com/DoodleScheduling/k8skafka-controller/kafka"
	"strconv"
)

const (
	CleanupPolicy                        = "cleanup.policy"
	CompressionType                      = "compression.type"
	DeleteRetentionMs                    = "delete.retention.ms"
	FileDeleteDelayMs                    = "file.delete.delay.ms"
	FlushMessages                        = "flush.messages"
	FlushMs                              = "flush.ms"
	FollowerReplicationThrottledReplicas = "follower.replication.throttled.replicas"
	IndexIntervalBytes                   = "index.interval.bytes"
	LeaderReplicationThrottledReplicas   = "leader.replication.throttled.replicas"
	MaxMessageBytes                      = "max.message.bytes"
	MessageDownconversionEnable          = "message.downconversion.enable"
	MessageFormatVersion                 = "message.format.version"
	MessageTimestampDifferenceMaxMs      = "message.timestamp.difference.max.ms"
	MessageTimestampType                 = "message.timestamp.type"
	MinCleanableDirtyRatio               = "min.cleanable.dirty.ratio"
	MinCompactionLagMs                   = "min.compaction.lag.ms"
	MinInsyncReplicas                    = "min.insync.replicas"
	Preallocate                          = "preallocate"
	RetentionBytes                       = "retention.bytes"
	RetentionMs                          = "retention.ms"
	SegmentBytes                         = "segment.bytes"
	SegmentIndexBytes                    = "segment.index.bytes"
	SegmentJitterMs                      = "segment.jitter.ms"
	SegmentMs                            = "segment.ms"
	UncleanLeaderElectionEnable          = "unclean.leader.election.enable"
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
		c[CleanupPolicy] = *topic.GetCleanupPolicy()
	}
	if topic.GetCompressionType() != nil {
		c[CompressionType] = *topic.GetCompressionType()
	}
	if topic.GetDeleteRetentionMs() != nil {
		c[DeleteRetentionMs] = strconv.FormatInt(*topic.GetDeleteRetentionMs(), 10)
	}
	if topic.GetFileDeleteDelayMs() != nil {
		c[FileDeleteDelayMs] = strconv.FormatInt(*topic.GetFileDeleteDelayMs(), 10)
	}
	if topic.GetFlushMessages() != nil {
		c[FlushMessages] = strconv.FormatInt(*topic.GetFlushMessages(), 10)
	}
	if topic.GetFlushMs() != nil {
		c[FlushMs] = strconv.FormatInt(*topic.GetFlushMs(), 10)
	}
	if topic.GetFollowerReplicationThrottledReplicas() != nil {
		c[FollowerReplicationThrottledReplicas] = *topic.GetFollowerReplicationThrottledReplicas()
	}
	if topic.GetIndexIntervalBytes() != nil {
		c[IndexIntervalBytes] = strconv.FormatInt(*topic.GetIndexIntervalBytes(), 10)
	}
	if topic.GetLeaderReplicationThrottledReplicas() != nil {
		c[LeaderReplicationThrottledReplicas] = *topic.GetLeaderReplicationThrottledReplicas()
	}
	if topic.GetMaxMessageBytes() != nil {
		c[MaxMessageBytes] = strconv.FormatInt(*topic.GetMaxMessageBytes(), 10)
	}
	if topic.GetMessageDownconversionEnable() != nil {
		c[MessageDownconversionEnable] = strconv.FormatBool(*topic.GetMessageDownconversionEnable())
	}
	if topic.GetMessageFormatVersion() != nil {
		c[MessageFormatVersion] = *topic.GetMessageFormatVersion()
	}
	if topic.GetMessageTimestampDifferenceMaxMs() != nil {
		c[MessageTimestampDifferenceMaxMs] = strconv.FormatInt(*topic.GetMessageTimestampDifferenceMaxMs(), 10)
	}
	if topic.GetMessageTimestampType() != nil {
		c[MessageTimestampType] = *topic.GetMessageTimestampType()
	}
	if topic.GetMinCleanableDirtyRatio() != nil {
		c[MinCleanableDirtyRatio] = strconv.FormatInt(*topic.GetMinCleanableDirtyRatio(), 10)
	}
	if topic.GetMinCompactionLagMs() != nil {
		c[MinCompactionLagMs] = strconv.FormatInt(*topic.GetMinCompactionLagMs(), 10)
	}
	if topic.GetMinInsyncReplicas() != nil {
		c[MinInsyncReplicas] = strconv.FormatInt(*topic.GetMinInsyncReplicas(), 10)
	}
	if topic.GetPreallocate() != nil {
		c[Preallocate] = strconv.FormatBool(*topic.GetPreallocate())
	}
	if topic.GetRetentionBytes() != nil {
		c[RetentionBytes] = strconv.FormatInt(*topic.GetRetentionBytes(), 10)
	}
	if topic.GetRetentionMs() != nil {
		c[RetentionMs] = strconv.FormatInt(*topic.GetRetentionMs(), 10)
	}
	if topic.GetSegmentBytes() != nil {
		c[SegmentBytes] = strconv.FormatInt(*topic.GetSegmentBytes(), 10)
	}
	if topic.GetSegmentIndexBytes() != nil {
		c[SegmentIndexBytes] = strconv.FormatInt(*topic.GetSegmentIndexBytes(), 10)
	}
	if topic.GetSegmentJitterMs() != nil {
		c[SegmentJitterMs] = strconv.FormatInt(*topic.GetSegmentJitterMs(), 10)
	}
	if topic.GetSegmentMs() != nil {
		c[SegmentMs] = strconv.FormatInt(*topic.GetSegmentMs(), 10)
	}
	if topic.GetUncleanLeaderElectionEnable() != nil {
		c[UncleanLeaderElectionEnable] = strconv.FormatBool(*topic.GetUncleanLeaderElectionEnable())
	}

	kt := kafka.Topic{
		Name:              topic.GetTopicName(),
		Partitions:        topic.GetPartitions(),
		ReplicationFactor: topic.GetReplicationFactor(),
		Config:            c,
	}

	return &kt
}
