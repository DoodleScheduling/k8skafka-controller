// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaTopic) DeepCopyInto(out *KafkaTopic) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaTopic.
func (in *KafkaTopic) DeepCopy() *KafkaTopic {
	if in == nil {
		return nil
	}
	out := new(KafkaTopic)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KafkaTopic) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaTopicConfig) DeepCopyInto(out *KafkaTopicConfig) {
	*out = *in
	if in.CleanupPolicy != nil {
		in, out := &in.CleanupPolicy, &out.CleanupPolicy
		*out = new(CleanupPolicy)
		**out = **in
	}
	if in.CompressionType != nil {
		in, out := &in.CompressionType, &out.CompressionType
		*out = new(CompressionType)
		**out = **in
	}
	if in.DeleteRetentionMs != nil {
		in, out := &in.DeleteRetentionMs, &out.DeleteRetentionMs
		*out = new(int64)
		**out = **in
	}
	if in.FileDeleteDelayMs != nil {
		in, out := &in.FileDeleteDelayMs, &out.FileDeleteDelayMs
		*out = new(int64)
		**out = **in
	}
	if in.FlushMessages != nil {
		in, out := &in.FlushMessages, &out.FlushMessages
		*out = new(int64)
		**out = **in
	}
	if in.FlushMs != nil {
		in, out := &in.FlushMs, &out.FlushMs
		*out = new(int64)
		**out = **in
	}
	if in.FollowerReplicationThrottledReplicas != nil {
		in, out := &in.FollowerReplicationThrottledReplicas, &out.FollowerReplicationThrottledReplicas
		*out = new(string)
		**out = **in
	}
	if in.IndexIntervalBytes != nil {
		in, out := &in.IndexIntervalBytes, &out.IndexIntervalBytes
		*out = new(int64)
		**out = **in
	}
	if in.LeaderReplicationThrottledReplicas != nil {
		in, out := &in.LeaderReplicationThrottledReplicas, &out.LeaderReplicationThrottledReplicas
		*out = new(string)
		**out = **in
	}
	if in.MaxMessageBytes != nil {
		in, out := &in.MaxMessageBytes, &out.MaxMessageBytes
		*out = new(int64)
		**out = **in
	}
	if in.MessageDownconversionEnable != nil {
		in, out := &in.MessageDownconversionEnable, &out.MessageDownconversionEnable
		*out = new(bool)
		**out = **in
	}
	if in.MessageFormatVersion != nil {
		in, out := &in.MessageFormatVersion, &out.MessageFormatVersion
		*out = new(string)
		**out = **in
	}
	if in.MessageTimestampDifferenceMaxMs != nil {
		in, out := &in.MessageTimestampDifferenceMaxMs, &out.MessageTimestampDifferenceMaxMs
		*out = new(int64)
		**out = **in
	}
	if in.MessageTimestampType != nil {
		in, out := &in.MessageTimestampType, &out.MessageTimestampType
		*out = new(MessageTimestampType)
		**out = **in
	}
	if in.MinCleanableDirtyRatio != nil {
		in, out := &in.MinCleanableDirtyRatio, &out.MinCleanableDirtyRatio
		*out = new(int64)
		**out = **in
	}
	if in.MinCompactionLagMs != nil {
		in, out := &in.MinCompactionLagMs, &out.MinCompactionLagMs
		*out = new(int64)
		**out = **in
	}
	if in.MinInsyncReplicas != nil {
		in, out := &in.MinInsyncReplicas, &out.MinInsyncReplicas
		*out = new(int64)
		**out = **in
	}
	if in.Preallocate != nil {
		in, out := &in.Preallocate, &out.Preallocate
		*out = new(bool)
		**out = **in
	}
	if in.RetentionBytes != nil {
		in, out := &in.RetentionBytes, &out.RetentionBytes
		*out = new(int64)
		**out = **in
	}
	if in.RetentionMs != nil {
		in, out := &in.RetentionMs, &out.RetentionMs
		*out = new(int64)
		**out = **in
	}
	if in.SegmentBytes != nil {
		in, out := &in.SegmentBytes, &out.SegmentBytes
		*out = new(int64)
		**out = **in
	}
	if in.SegmentIndexBytes != nil {
		in, out := &in.SegmentIndexBytes, &out.SegmentIndexBytes
		*out = new(int64)
		**out = **in
	}
	if in.SegmentJitterMs != nil {
		in, out := &in.SegmentJitterMs, &out.SegmentJitterMs
		*out = new(int64)
		**out = **in
	}
	if in.SegmentMs != nil {
		in, out := &in.SegmentMs, &out.SegmentMs
		*out = new(int64)
		**out = **in
	}
	if in.UncleanLeaderElectionEnable != nil {
		in, out := &in.UncleanLeaderElectionEnable, &out.UncleanLeaderElectionEnable
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaTopicConfig.
func (in *KafkaTopicConfig) DeepCopy() *KafkaTopicConfig {
	if in == nil {
		return nil
	}
	out := new(KafkaTopicConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaTopicList) DeepCopyInto(out *KafkaTopicList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KafkaTopic, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaTopicList.
func (in *KafkaTopicList) DeepCopy() *KafkaTopicList {
	if in == nil {
		return nil
	}
	out := new(KafkaTopicList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KafkaTopicList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaTopicSpec) DeepCopyInto(out *KafkaTopicSpec) {
	*out = *in
	if in.Partitions != nil {
		in, out := &in.Partitions, &out.Partitions
		*out = new(int64)
		**out = **in
	}
	if in.ReplicationFactor != nil {
		in, out := &in.ReplicationFactor, &out.ReplicationFactor
		*out = new(int64)
		**out = **in
	}
	if in.KafkaTopicConfig != nil {
		in, out := &in.KafkaTopicConfig, &out.KafkaTopicConfig
		*out = new(KafkaTopicConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaTopicSpec.
func (in *KafkaTopicSpec) DeepCopy() *KafkaTopicSpec {
	if in == nil {
		return nil
	}
	out := new(KafkaTopicSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaTopicStatus) DeepCopyInto(out *KafkaTopicStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaTopicStatus.
func (in *KafkaTopicStatus) DeepCopy() *KafkaTopicStatus {
	if in == nil {
		return nil
	}
	out := new(KafkaTopicStatus)
	in.DeepCopyInto(out)
	return out
}