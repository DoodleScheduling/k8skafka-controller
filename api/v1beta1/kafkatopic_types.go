/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultPartitions        = 1
	defaultReplicationFactor = 1
)

// KafkaTopicSpec defines the desired state of KafkaTopic
type KafkaTopicSpec struct {

	// The connect URI
	// +required
	Address string `json:"address"`

	// Name is by default the same as metata.name
	// +optional
	Name string `json:"name"`

	// Number of partitions
	// +optional
	Partitions *int64 `json:"partitions,omitempty"`

	// Replication factor
	// +optional
	ReplicationFactor *int64 `json:"replicationFactor,omitempty"`

	// Additional topic configuration
	// +optional
	KafkaTopicConfig *KafkaTopicConfig `json:"config,omitempty"`
}

// KafkaTopicConfig defines additional topic configuration
type KafkaTopicConfig struct {

	// Designates the retention policy to use on old log segments.
	// Either "delete" or "compact" or both ("delete,compact").
	// The default policy ("delete") will discard old segments when their retention time or size limit has been reached.
	// The "compact" setting will enable log compaction on the topic.
	// +optional
	CleanupPolicy *CleanupPolicy `json:"cleanupPolicy,omitempty"`

	// Final compression type for a given topic.
	// Supported are standard compression codecs: 'gzip', 'snappy', 'lz4', 'zstd').
	// It additionally accepts 'uncompressed' which is equivalent to no compression; and 'producer' which means retain the original compression codec set by the producer.
	// +optional
	CompressionType *CompressionType `json:"compressionType,omitempty"`

	// The amount of time to retain delete tombstone markers for log compacted topics. Specified in milliseconds.
	// This setting also gives a bound on the time in which a consumer must complete a read if they begin from offset 0
	// to ensure that they get a valid snapshot of the final stage (otherwise delete tombstones may be collected before they complete their scan).
	// +optional
	DeleteRetentionMs *int64 `json:"deleteRetentionsMs,omitempty"`

	// The time to wait before deleting a file from the filesystems
	// +optional
	FileDeleteDelayMs *int64 `json:"fileDeleteDelayMs,omitempty"`

	// This setting allows specifying an interval at which there will be a force if an fsync of data written to the log.
	// For example, if this was set to 1 there would be a fsync after every message; if it were 5 there would be a fsync after every five messages.
	// In general, it is recommended not to set this and use replication for durability and allow the operating system's background flush capabilities as it is more efficient.
	// +optional
	FlushMessages *int64 `json:"flushMessages,omitempty"`

	// This setting allows specifying a time interval at which there will be a force of an fsync of data written to the log.
	// For example if this was set to 1000 there would be a fsync after 1000 ms had passed.
	// In general, it is not recommended to set this and instead use replication for durability and allow the operating system's background flush capabilities as it is more efficient.
	// +optional
	FlushMs *int64 `json:"flushMs,omitempty"`

	// A list of replicas for which log replication should be throttled on the follower side.
	// The list should describe a set of replicas in the form [PartitionId]:[BrokerId],[PartitionId]:[BrokerId]:... or alternatively the wildcard '*' can be used to throttle all replicas for this topic.
	// +optional
	FollowerReplicationThrottledReplicas *string `json:"followerReplicationThrottledReplicas,omitempty"`

	// This setting controls how frequently Kafka adds an index entry to its offset index.
	// The default setting ensures that a messages is indexed roughly every 4096 bytes.
	// More indexing allows reads to jump closer to the exact position in the log but makes the index larger. You probably don't need to change this.
	// +optional
	IndexIntervalBytes *int64 `json:"indexIntervalBytes,omitempty"`

	// A list of replicas for which log replication should be throttled on the leader side.
	// The list should describe a set of replicas in the form [PartitionId]:[BrokerId],[PartitionId]:[BrokerId]:... or alternatively the wildcard '*' can be used to throttle all replicas for this topic.
	// +optional
	LeaderReplicationThrottledReplicas *string `json:"leaderReplicationThrottledReplicas,omitempty"`

	// The largest record batch size allowed by Kafka.
	// If this is increased and there are consumers older than 0.10.2, the consumers' fetch size must also be increased so that the they can fetch record batches this large.
	// In the latest message format version, records are always grouped into batches for efficiency.
	// In previous message format versions, uncompressed records are not grouped into batches and this limit only applies to a single record in that case.
	// +optional
	MaxMessageBytes *int64 `json:"maxMessageBytes,omitempty"`

	// This configuration controls whether down-conversion of message formats is enabled to satisfy consume requests.
	// When set to false, broker will not perform down-conversion for consumers expecting an older message format.
	// The broker responds with UNSUPPORTED_VERSION error for consume requests from such older clients.
	// This configuration does not apply to any message format conversion that might be required for replication to followers.
	// +optional
	MessageDownconversionEnable *bool `json:"messageDownconversionEnable,omitempty"`

	// Specify the message format version the broker will use to append messages to the logs.
	// The value should be a valid ApiVersion. Some examples are: 0.8.2, 0.9.0.0, 0.10.0, check ApiVersion for more details.
	// By setting a particular message format version, the user is certifying that all the existing messages on disk are smaller or equal than the specified version.
	// Setting this value incorrectly will cause consumers with older versions to break as they will receive messages with a format that they don't understand.
	// +optional
	MessageFormatVersion *string `json:"messageFormatVersion,omitempty"`

	// The maximum difference allowed between the timestamp when a broker receives a message and the timestamp specified in the message.
	// If MessageTimestampType=CreateTime, a message will be rejected if the difference in timestamp exceeds this threshold.
	// This configuration is ignored if MessageTimestampType=LogAppendTime.
	// +optional
	MessageTimestampDifferenceMaxMs *int64 `json:"messageTimestampDifferenceMaxMs,omitempty"`

	// Define whether the timestamp in the message is message create time or log append time.
	// The value should be either `CreateTime` or `LogAppendTime`
	// +optional
	MessageTimestampType *MessageTimestampType `json:"messageTimestampType,omitempty"`

	// This configuration controls how frequently the log compactor will attempt to clean the log (assuming LogCompaction is enabled).
	// By default we will avoid cleaning a log where more than 50% of the log has been compacted.
	// This ratio bounds the maximum space wasted in the log by duplicates (at 50% at most 50% of the log could be duplicates).
	// A higher ratio will mean fewer, more efficient cleanings but will mean more wasted space in the log.
	// If the MaxCompactionLagMs or the MinCompactionLagMs configurations are also specified, then the log compactor considers the log to be eligible for compaction as soon as either:
	// (i) the dirty ratio threshold has been met and the log has had dirty (uncompacted) records for at least the MinCompactionLagMs duration,
	// or (ii) if the log has had dirty (uncompacted) records for at most the MaxCompactionLagMs period.
	// +optional
	MinCleanableDirtyRatio *int64 `json:"minCleanableDirtyRatio,omitempty"`

	// The minimum time a message will remain uncompacted in the log. Only applicable for logs that are being compacted.
	// +optional
	MinCompactionLagMs *int64 `json:"minCompactionLagMs,omitempty"`

	// When a producer sets acks to "all" (or "-1"), this configuration specifies the minimum number of replicas that must acknowledge a write for the write to be considered successful.
	// If this minimum cannot be met, then the producer will raise an exception (either NotEnoughReplicas or NotEnoughReplicasAfterAppend).
	// When used together, MinInsyncReplicas and acks allow you to enforce greater durability guarantees.
	// A typical scenario would be to create a topic with a replication factor of 3, set MinInsyncReplicas to 2, and produce with ack of "all".
	// This will ensure that the producer raises an exception if a majority of replicas do not receive a write.
	// +optional
	MinInsyncReplicas *int64 `json:"minInsyncReplicas,omitempty"`

	// True if we should preallocate the file on disk when creating a new log segment.
	// +optional
	Preallocate *bool `json:"preallocate,omitempty"`

	// This configuration controls the maximum size a partition (which consists of log segments) can grow to before we will discard old log segments to free up space if we are using the "delete" retention policy.
	// By default there is no size limit only a time limit. Since this limit is enforced at the partition level, multiply it by the number of partitions to compute the topic retention in bytes.
	// +optional
	RetentionBytes *int64 `json:"retentionBytes,omitempty"`

	// This configuration controls the maximum time we will retain a log before we will discard old log segments to free up space if we are using the "delete" retention policy.
	// This represents an SLA on how soon consumers must read their data. If set to -1, no time limit is applied.
	// +optional
	RetentionMs *int64 `json:"retentionMs,omitempty"`

	// This configuration controls the segment file size for the log. Retention and cleaning is always done a file at a time so a larger segment size means fewer files but less granular control over retention.
	// +optional
	SegmentBytes *int64 `json:"segmentBytes,omitempty"`

	// This configuration controls the size of the index that maps offsets to file positions. We preallocate this index file and shrink it only after log rolls.
	// You generally should not need to change this setting.
	// +optional
	SegmentIndexBytes *int64 `json:"segmentIndexBytes,omitempty"`

	// The maximum random jitter subtracted from the scheduled segment roll time to avoid thundering herds of segment rolling
	// +optional
	SegmentJitterMs *int64 `json:"segmentJitterMs,omitempty"`

	// This configuration controls the period of time after which Kafka will force the log to roll even if the segment file isn't full to ensure that retention can delete or compact old data.
	// +optional
	SegmentMs *int64 `json:"segmentMs,omitempty"`

	// Indicates whether to enable replicas not in the ISR set to be elected as leader as a last resort, even though doing so may result in data loss.
	// +optional
	UncleanLeaderElectionEnable *bool `json:"uncleanLeaderElectionEnable,omitempty"`
}

type CleanupPolicy string

const (
	CleanupPolicyDelete        = "delete"
	CleanupPolicyCompact       = "compact"
	CleanupPolicyDeleteCompact = "delete,compact"
)

type CompressionType string

const (
	CompressionTypeGZIP         = "gzip"
	CompressionTypeSnappy       = "snappy"
	CompressionTypeLZ4          = "lz4"
	CompressionTypeZSTD         = "zstd"
	CompressionTypeUncompressed = "uncompressed"
	CompressionTypeProducer     = "producer"
)

type MessageTimestampType string

const (
	MessageTimestampTypeCreateTime    = "CreateTime"
	MessageTimestampTypeLogAppendTime = "LogAppendTime"
)

// KafkaTopicStatus defines the observed state of KafkaTopic
type KafkaTopicStatus struct {
	// Conditions holds the conditions for the KafkaTopic.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

const (
	ReadyCondition = "Ready"
)

const (
	ReplicationFactorFailedToChangeReason = "ReplicationFactorFailedToChange"
	PartitionsFailedToRemoveReason        = "PartitionsFailedToRemove"
	PartitionsFailedToCreateReason        = "PartitionsFailedToCreate"
	TopicFailedToGetReason                = "TopicFailedToGet"
	TopicFailedToCreateReason             = "TopicFailedToCreate"
	TopicReadyReason                      = "TopicReadyReason"
)

// ConditionalResource is a resource with conditions
type conditionalResource interface {
	GetStatusConditions() *[]metav1.Condition
}

// setResourceCondition sets the given condition with the given status,
// reason and message on a resource.
func setResourceCondition(resource conditionalResource, condition string, status metav1.ConditionStatus, reason, message string) {
	conditions := resource.GetStatusConditions()

	newCondition := metav1.Condition{
		Type:    condition,
		Status:  status,
		Reason:  reason,
		Message: message,
	}

	apimeta.SetStatusCondition(conditions, newCondition)
}

// KafkaTopicReady
func KafkaTopicReady(topic KafkaTopic, reason, message string) KafkaTopic {
	setResourceCondition(&topic, ReadyCondition, metav1.ConditionTrue, reason, message)
	return topic
}

// KafkaTopicNotReady
func KafkaTopicNotReady(topic KafkaTopic, reason, message string) KafkaTopic {
	setResourceCondition(&topic, ReadyCondition, metav1.ConditionFalse, reason, message)
	return topic
}

// GetStatusConditions returns a pointer to the Status.Conditions slice
func (in *KafkaTopic) GetStatusConditions() *[]metav1.Condition {
	return &in.Status.Conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=kt
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description=""
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description=""

// KafkaTopic is the Schema for the kafkatopics API
type KafkaTopic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KafkaTopicSpec   `json:"spec,omitempty"`
	Status KafkaTopicStatus `json:"status,omitempty"`
}

func (in *KafkaTopic) GetAddress() string {
	return in.Spec.Address
}

func (in *KafkaTopic) GetTopicName() string {
	if in.Spec.Name != "" {
		return in.Spec.Name
	}
	return in.GetName()
}

func (in *KafkaTopic) GetPartitions() int64 {
	if in.Spec.Partitions == nil {
		return defaultPartitions
	}
	return *in.Spec.Partitions
}

func (in *KafkaTopic) GetReplicationFactor() int64 {
	if in.Spec.ReplicationFactor == nil {
		return defaultReplicationFactor
	}
	return *in.Spec.ReplicationFactor
}

func (in *KafkaTopic) GetCleanupPolicy() *CleanupPolicy {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.CleanupPolicy
}

func (in *KafkaTopic) GetCompressionType() *CompressionType {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.CompressionType
}

func (in *KafkaTopic) GetDeleteRetentionMs() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.DeleteRetentionMs
}

func (in *KafkaTopic) GetFileDeleteDelayMs() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.FileDeleteDelayMs
}

func (in *KafkaTopic) GetFlushMessages() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.FlushMessages
}

func (in *KafkaTopic) GetFlushMs() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.FlushMs
}

func (in *KafkaTopic) GetFollowerReplicationThrottledReplicas() *string {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.FollowerReplicationThrottledReplicas
}

func (in *KafkaTopic) GetIndexIntervalBytes() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.IndexIntervalBytes
}

func (in *KafkaTopic) GetLeaderReplicationThrottledReplicas() *string {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.LeaderReplicationThrottledReplicas
}

func (in *KafkaTopic) GetMaxMessageBytes() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.MaxMessageBytes
}

func (in *KafkaTopic) GetMessageDownconversionEnable() *bool {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.MessageDownconversionEnable
}

func (in *KafkaTopic) GetMessageFormatVersion() *string {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.MessageFormatVersion
}

func (in *KafkaTopic) GetMessageTimestampDifferenceMaxMs() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.MessageTimestampDifferenceMaxMs
}

func (in *KafkaTopic) GetMessageTimestampType() *MessageTimestampType {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.MessageTimestampType
}

func (in *KafkaTopic) GetMinCleanableDirtyRatio() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.MinCleanableDirtyRatio
}

func (in *KafkaTopic) GetMinCompactionLagMs() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.MinCompactionLagMs
}

func (in *KafkaTopic) GetMinInsyncReplicas() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.MinInsyncReplicas
}

func (in *KafkaTopic) GetPreallocate() *bool {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.Preallocate
}

func (in *KafkaTopic) GetRetentionBytes() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.RetentionBytes
}

func (in *KafkaTopic) GetRetentionMs() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.RetentionMs
}

func (in *KafkaTopic) GetSegmentBytes() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.SegmentBytes
}

func (in *KafkaTopic) GetSegmentIndexBytes() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.SegmentIndexBytes
}

func (in *KafkaTopic) GetSegmentJitterMs() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.SegmentJitterMs
}

func (in *KafkaTopic) GetSegmentMs() *int64 {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.SegmentMs
}

func (in *KafkaTopic) GetUncleanLeaderElectionEnable() *bool {
	if in.Spec.KafkaTopicConfig == nil {
		return nil
	}
	return in.Spec.KafkaTopicConfig.UncleanLeaderElectionEnable
}

// +kubebuilder:object:root=true

// KafkaTopicList contains a list of KafkaTopic
type KafkaTopicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KafkaTopic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KafkaTopic{}, &KafkaTopicList{})
}
