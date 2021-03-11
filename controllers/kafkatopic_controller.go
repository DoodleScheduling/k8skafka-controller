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
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

type KafkaTopicReconcilerOptions struct {
	MaxConcurrentReconciles int
}

// +kubebuilder:rbac:groups=kafka.infra.doodle.com,resources=KafkaTopics,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kafka.infra.doodle.com,resources=KafkaTopics/status,verbs=get;update;patch

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

	topic, result, reconcileErr := r.reconcile(ctx, topic)

	// Update status after reconciliation.
	if err = r.patchStatus(ctx, &topic); err != nil {
		logger.Error(err, "unable to update status after reconciliation")
		return ctrl.Result{Requeue: true}, err
	}

	return result, reconcileErr
}

func (r *KafkaTopicReconciler) reconcile(ctx context.Context, topic v1beta1.KafkaTopic) (v1beta1.KafkaTopic, ctrl.Result, error) {
	kc := kafka.NewTCPConnection(topic.GetAddress())
	if err := kc.CreatePartitions(ctx, kafka.Topic{
		Name:              topic.GetTopicName(),
		Partitions:        topic.GetPartitions(),
		ReplicationFactor: topic.GetReplicationFactor(),
	}); err != nil {
		msg := fmt.Sprintf("Topic/partitions failed to create: %s", err.Error())
		r.Recorder.Event(&topic, "Normal", "info", msg)
		return v1beta1.KafkaTopicNotReady(topic, v1beta1.TopicFailedToCreateReason, msg), ctrl.Result{}, err
	} else {
		msg := "Topic/partitions successfully created"
		r.Recorder.Event(&topic, "Normal", "info", msg)
		return v1beta1.KafkaTopicReady(topic, v1beta1.TopicReadyReason, msg), ctrl.Result{}, err
	}
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
