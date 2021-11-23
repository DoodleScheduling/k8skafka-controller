package controllers

import infrav1beta1 "github.com/DoodleScheduling/k8skafka-controller/api/v1beta1"

// KafkaTopicConfigHolder
/* KafkaTopicConfigHolder holds testing data.
valueToSet : new value to set for config option in test
expectedValue : value that is expected as result of the test
createKafkaTopicObjectF : function that creates KafkaTopic k8s object from given config value input
*/
type KafkaTopicConfigHolder struct {
	valueToSet              interface{}
	expectedValue           interface{}
	createKafkaTopicObjectF func(interface{}) *infrav1beta1.KafkaTopicConfig
}

// KafkaTopicConfigTestData
/* KafkaTopicConfigTestData holds testing data for each topic configuration option.
It is a map of slices, where each entry in the map represents one configuration option, and entries in slices are possible tests variations.
To add a new test for existing configuration option, simply create new entry in slice for appropriate config option
To add a completely new config option, add a new key with map (value is a slice of KafkaTopicConfigHolder)
To change existing tests, simply change it :)
*/
var KafkaTopicConfigTestData = map[string][]KafkaTopicConfigHolder{
	CleanupPolicy: {
		{
			infrav1beta1.CleanupPolicyCompact, infrav1beta1.CleanupPolicyCompact, func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(string)
				return &infrav1beta1.KafkaTopicConfig{
					CleanupPolicy: &vv,
				}
			},
		},
	},
	CompressionType: {
		{
			infrav1beta1.CompressionTypeSnappy, infrav1beta1.CompressionTypeSnappy, func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(string)
				return &infrav1beta1.KafkaTopicConfig{
					CompressionType: &vv,
				}
			},
		},
	},
	DeleteRetentionMs: {
		{
			int64(60000), "60000", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					DeleteRetentionMs: &vv,
				}
			},
		},
	},
	FileDeleteDelayMs: {
		{
			int64(1000), "1000", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					FileDeleteDelayMs: &vv,
				}
			},
		},
	},
	FlushMessages: {
		{
			int64(5), "5", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					FlushMessages: &vv,
				}
			},
		},
	},
	FlushMs: {
		{
			int64(666), "666", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					FlushMs: &vv,
				}
			},
		},
	},
	FollowerReplicationThrottledReplicas: {
		{
			"*", "*", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(string)
				return &infrav1beta1.KafkaTopicConfig{
					FollowerReplicationThrottledReplicas: &vv,
				}
			},
		},
	},
	IndexIntervalBytes: {
		{
			int64(2048), "2048", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					IndexIntervalBytes: &vv,
				}
			},
		},
	},
	LeaderReplicationThrottledReplicas: {
		{
			"*", "*", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(string)
				return &infrav1beta1.KafkaTopicConfig{
					LeaderReplicationThrottledReplicas: &vv,
				}
			},
		},
	},
	MaxMessageBytes: {
		{
			int64(999999), "999999", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					MaxMessageBytes: &vv,
				}
			},
		},
	},
	MessageDownconversionEnable: {
		{
			true, "true", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(bool)
				return &infrav1beta1.KafkaTopicConfig{
					MessageDownconversionEnable: &vv,
				}
			},
		},
	},
	MessageFormatVersion: {
		{
			"0.10.0", "0.10.0", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(string)
				return &infrav1beta1.KafkaTopicConfig{
					MessageFormatVersion: &vv,
				}
			},
		},
	},
	MessageTimestampDifferenceMaxMs: {
		{
			int64(10), "10", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					MessageTimestampDifferenceMaxMs: &vv,
				}
			},
		},
	},
	MessageTimestampType: {
		{
			"LogAppendTime", "LogAppendTime", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(string)
				return &infrav1beta1.KafkaTopicConfig{
					MessageTimestampType: &vv,
				}
			},
		},
	},
	MinCleanableDirtyRatio: {
		{
			int64(50), "50", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					MinCleanableDirtyRatio: &vv,
				}
			},
		},
	},
	MinCompactionLagMs: {
		{
			int64(10000), "10000", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					MinCompactionLagMs: &vv,
				}
			},
		},
	},
	MinInsyncReplicas: {
		{
			int64(2), "2", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					MinInsyncReplicas: &vv,
				}
			},
		},
	},
	Preallocate: {
		{
			true, "true", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(bool)
				return &infrav1beta1.KafkaTopicConfig{
					Preallocate: &vv,
				}
			},
		},
	},
	RetentionBytes: {
		{
			int64(1000000), "1000000", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					RetentionBytes: &vv,
				}
			},
		},
	},
	RetentionMs: {
		{
			int64(1000), "1000", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					RetentionMs: &vv,
				}
			},
		},
	},
	SegmentBytes: {
		{
			int64(500000), "500000", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					SegmentBytes: &vv,
				}
			},
		},
	},
	SegmentIndexBytes: {
		{
			int64(250000), "250000", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					SegmentIndexBytes: &vv,
				}
			},
		},
	},
	SegmentJitterMs: {
		{
			int64(1000), "1000", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					SegmentJitterMs: &vv,
				}
			},
		},
	},
	SegmentMs: {
		{
			int64(2000), "2000", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(int64)
				return &infrav1beta1.KafkaTopicConfig{
					SegmentMs: &vv,
				}
			},
		},
	},
	UncleanLeaderElectionEnable: {
		{
			true, "true", func(v interface{}) *infrav1beta1.KafkaTopicConfig {
				vv := v.(bool)
				return &infrav1beta1.KafkaTopicConfig{
					UncleanLeaderElectionEnable: &vv,
				}
			},
		},
	},
}
