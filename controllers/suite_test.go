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
	"github.com/DoodleScheduling/k8skafka-controller/kafka"
	k "github.com/segmentio/kafka-go"
	"net"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrav1beta1 "github.com/DoodleScheduling/k8skafka-controller/api/v1beta1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg                     *rest.Config
	k8sClient               client.Client
	testEnv                 *envtest.Environment
	ctx                     context.Context
	cancel                  context.CancelFunc
	kafkaCluster            *TestingKafkaCluster
	TestingKafkaClusterHost string
	//TopicValidator
)

const (
	numberOfConcurrentReconcilers = 1
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

func GetTopic(name string) (*k.Topic, error) {
	addr, err := net.ResolveTCPAddr("tcp", TestingKafkaClusterHost)
	if err != nil {
		return nil, err
	}
	kc := k.Client{
		Addr:    addr,
		Timeout: 4 * time.Minute,
	}
	metadataResponse, err := kc.Metadata(context.Background(), &k.MetadataRequest{
		Addr:   addr,
		Topics: []string{name},
	})
	if err != nil {
		return nil, err
	}
	for _, topic := range metadataResponse.Topics {
		if topic.Name == name {
			if topic.Error != nil {
				return nil, topic.Error
			}
			return &topic, nil
		}
	}
	return nil, nil
}

func GetTopicConfig(name string) ([]k.DescribeConfigResponseResource, error) {
	addr, err := net.ResolveTCPAddr("tcp", TestingKafkaClusterHost)
	if err != nil {
		return nil, err
	}
	kc := k.Client{
		Addr:    addr,
		Timeout: 4 * time.Minute,
	}
	describeConfigsResponse, err := kc.DescribeConfigs(context.Background(), &k.DescribeConfigsRequest{
		Addr: addr,
		Resources: []k.DescribeConfigRequestResource{
			{
				ResourceType: k.ResourceTypeTopic,
				ResourceName: name,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return describeConfigsResponse.Resources, nil
}

var _ = BeforeSuite(func(done Done) {
	//logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	By("setting up kafka cluster")
	kafkaCluster, err = NewTestingKafkaCluster()
	Expect(err).NotTo(HaveOccurred())
	Expect(kafkaCluster).ToNot(BeNil())

	By("starting kafka cluster")
	err = kafkaCluster.StartCluster()
	Expect(err).ToNot(HaveOccurred())

	TestingKafkaClusterHost, err = kafkaCluster.GetKafkaHost()
	Expect(err).ToNot(HaveOccurred())
	Expect(TestingKafkaClusterHost).ToNot(BeEmpty())

	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = infrav1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = infrav1beta1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	opts := KafkaTopicReconcilerOptions{MaxConcurrentReconciles: numberOfConcurrentReconcilers}

	err = (&KafkaTopicReconciler{
		Client:      k8sManager.GetClient(),
		Scheme:      k8sManager.GetScheme(),
		Log:         logf.Log,
		Recorder:    k8sManager.GetEventRecorderFor("KafkaTopic"),
		KafkaClient: kafka.NewDefaultKafkaClient(),
	}).SetupWithManager(k8sManager, opts)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())

	err = kafkaCluster.StopCluster()
	Expect(err).ToNot(HaveOccurred())
})
