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
}

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
