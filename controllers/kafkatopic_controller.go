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

package controllers

import (
	"context"
	"fmt"
	"github.com/DoodleScheduling/k8skafka-controller/kafka"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1beta1 "github.com/DoodleScheduling/k8skafka-controller/api/v1beta1"
)

// KafkaTopicReconciler reconciles a KafkaTopic object
type KafkaTopicReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	Recorder    record.EventRecorder
	KafkaClient kafka.KafkaClient
}

type KafkaTopicReconcilerOptions struct {
	MaxConcurrentReconciles int
}

// +kubebuilder:rbac:groups=kafka.infra.doodle.com,resources=kafkatopics,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kafka.infra.doodle.com,resources=kafkatopics/status,verbs=get;update;patch

func (r *KafkaTopicReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("Namespace", req.Namespace, "Name", req.NamespacedName)
	logger.Info("reconciling KafkaTopic")

	// Fetch the RequestClone instance
	topic := v1beta1.KafkaTopic{}

	err := r.Client.Get(ctx, req.NamespacedName, &topic)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	topic, result := r.reconcile(ctx, topic)

	// Update status after reconciliation.
	if err = r.patchStatus(ctx, &topic); err != nil {
		logger.Error(err, "unable to update status after reconciliation")
		return ctrl.Result{Requeue: true}, err
	}

	return result, nil
}

func (r *KafkaTopicReconciler) reconcile(ctx context.Context, topic v1beta1.KafkaTopic) (v1beta1.KafkaTopic, ctrl.Result) {
	kt := TranslateKafkaTopicV1Beta1(topic)

	existingTopic, err := r.KafkaClient.GetTopic(topic.GetAddress(), kt.Name)
	if err != nil {
		msg := fmt.Sprintf("Cannot get topic: %s in %s :: %s", kt.Name, topic.GetAddress(), err.Error())
		r.Recorder.Event(&topic, "Normal", "info", msg)
		return v1beta1.KafkaTopicNotReady(topic, v1beta1.TopicFailedToGetReason, msg), ctrl.Result{Requeue: true}
	}
	if existingTopic == nil {
		if err := r.KafkaClient.CreateTopic(topic.GetAddress(), *kt); err != nil {
			msg := fmt.Sprintf("Failed to create topic: %s", err.Error())
			r.Recorder.Event(&topic, "Normal", "info", msg)
			return v1beta1.KafkaTopicNotReady(topic, v1beta1.TopicFailedToCreateReason, msg), ctrl.Result{Requeue: true}
		}

		msg := "Topic successfully created."
		r.Recorder.Event(&topic, "Normal", "info", msg)
		return v1beta1.KafkaTopicReady(topic, v1beta1.TopicReadyReason, msg), ctrl.Result{}
	}

	if existingTopic.ReplicationFactor != kt.ReplicationFactor {
		msg := fmt.Sprintf("Cannot change replication factor, this is currently not supported. "+
			"Requested replication factor: %d, current replication factor: %d, topic: %s, address: %s", kt.ReplicationFactor, existingTopic.ReplicationFactor, kt.Name, topic.GetAddress())
		r.Recorder.Event(&topic, "Normal", "info", msg)
		return v1beta1.KafkaTopicNotReady(topic, v1beta1.ReplicationFactorFailedToChangeReason, msg), ctrl.Result{}
	}
	if existingTopic.Partitions != kt.Partitions {
		if existingTopic.Partitions > kt.Partitions {
			msg := fmt.Sprintf("Cannot remove partitions, this is not allowed. "+
				"Requested number of partitions: %d, current partitions: %d, topic: %s, address: %s", kt.Partitions, existingTopic.Partitions, kt.Name, topic.GetAddress())
			r.Recorder.Event(&topic, "Normal", "info", msg)
			return v1beta1.KafkaTopicNotReady(topic, v1beta1.PartitionsFailedToRemoveReason, msg), ctrl.Result{}
		}

		kt.Brokers = existingTopic.Brokers
		if err := r.KafkaClient.CreatePartitions(ctx, topic.GetAddress(), *kt, kt.Partitions-existingTopic.Partitions); err != nil {
			msg := fmt.Sprintf("Failed to create partitions: %s", err.Error())
			r.Recorder.Event(&topic, "Normal", "info", msg)
			return v1beta1.KafkaTopicNotReady(topic, v1beta1.PartitionsFailedToCreateReason, msg), ctrl.Result{Requeue: true}
		}
	}
	if err := r.KafkaClient.UpdateTopicConfiguration(ctx, topic.GetAddress(), *kt); err != nil {
		msg := fmt.Sprintf("Failed to update topic: %s %s", kt.Name, err.Error())
		r.Recorder.Event(&topic, "Normal", "info", msg)
		return v1beta1.KafkaTopicNotReady(topic, v1beta1.TopicFailedToUpdateReason, msg), ctrl.Result{Requeue: true}
	}
	msg := "Topic successfully updated."
	r.Recorder.Event(&topic, "Normal", "info", msg)
	return v1beta1.KafkaTopicReady(topic, v1beta1.TopicReadyReason, msg), ctrl.Result{}
}

func (r *KafkaTopicReconciler) SetupWithManager(mgr ctrl.Manager, opts KafkaTopicReconcilerOptions) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.KafkaTopic{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: opts.MaxConcurrentReconciles}).
		Complete(r)
}

func (r *KafkaTopicReconciler) patchStatus(ctx context.Context, topic *v1beta1.KafkaTopic) error {
	key := client.ObjectKeyFromObject(topic)
	latest := &v1beta1.KafkaTopic{}
	if err := r.Client.Get(ctx, key, latest); err != nil {
		return err
	}

	return r.Client.Status().Patch(ctx, topic, client.MergeFrom(latest))
}
