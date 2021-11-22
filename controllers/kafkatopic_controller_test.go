package controllers

import (
	"context"
	"github.com/DoodleScheduling/k8skafka-controller/kafka"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	infrav1beta1 "github.com/DoodleScheduling/k8skafka-controller/api/v1beta1"
)

var _ = Describe("KafkaTopic controller", func() {
	const (
		KafkaBrokersAddress = "kafka-client.default:9092"
		KafkaTopicName      = "test-kafka-topic"
		KafkaTopicNamespace = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	BeforeEach(func() {
		kafka.DefaultMockKafkaBrokers.ClearAllTopics()
	})

	AfterEach(func() {

	})

	Context("When creating new topic that doesn't exist already", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 3

		It("Should create new topic", func() {
			By("By creating a KafkaTopic object in API server")
			ctx := context.Background()
			kafkaTopic := &infrav1beta1.KafkaTopic{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "kafka.infra.doodle.com/v1beta1",
					Kind:       "KafkaTopic",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KafkaTopicName,
					Namespace: KafkaTopicNamespace,
				},
				Spec: infrav1beta1.KafkaTopicSpec{
					Address:           KafkaBrokersAddress,
					Name:              KafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())

			kafkaTopicLookupKey := types.NamespacedName{Name: KafkaTopicName, Namespace: KafkaTopicNamespace}
			createdKafkaTopic := &infrav1beta1.KafkaTopic{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, createdKafkaTopic)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdKafkaTopic.Spec.Name).Should(Equal(KafkaTopicName))

			By("By checking that partitions are as expected")
			Expect(*createdKafkaTopic.Spec.Partitions).Should(Equal(partitions))

			By("By checking that replication factor is as expected")
			Expect(*createdKafkaTopic.Spec.ReplicationFactor).Should(Equal(replicationFactor))

			By("By checking that kafka brokers address is as expected")
			Expect(createdKafkaTopic.Spec.Address).Should(Equal(KafkaBrokersAddress))

			By("By creating a kafka topic in brokers")
			Eventually(func() bool {
				return kafka.DefaultMockKafkaBrokers.GetTopic(KafkaTopicName) != nil
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				topic := kafka.DefaultMockKafkaBrokers.GetTopic(KafkaTopicName)
				if topic == nil {
					return false
				}
				return topic.Name == KafkaTopicName
			}, timeout, interval).Should(BeTrue())

			By("By checking that kafka topic in brokers has expected number of partitions")
			Eventually(func() bool {
				topic := kafka.DefaultMockKafkaBrokers.GetTopic(KafkaTopicName)
				if topic == nil {
					return false
				}
				return topic.Partitions == partitions
			}).Should(BeTrue())

			By("By checking that kafka topic in brokers has expected replication factor")
			Eventually(func() bool {
				topic := kafka.DefaultMockKafkaBrokers.GetTopic(KafkaTopicName)
				if topic == nil {
					return false
				}
				return topic.ReplicationFactor == replicationFactor
			}).Should(BeTrue())
		})
	})
})
