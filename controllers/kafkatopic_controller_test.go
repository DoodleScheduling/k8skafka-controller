package controllers

import (
	"context"
	"github.com/DoodleScheduling/k8skafka-controller/kafka"
	"github.com/pkg/errors"
	"strconv"
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
		KafkaTopicNamespace = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	BeforeEach(func() {
		//kafka.DefaultMockKafkaBrokers.ClearAllTopics()
	})

	AfterEach(func() {
	})

	Context("When creating new topic that doesn't exist already", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 3
		kafkaTopicName := "test-new"

		It("Should create new topic", func() {
			By("By creating a KafkaTopic object in API server")
			ctx := context.Background()
			kafkaTopic := &infrav1beta1.KafkaTopic{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "kafka.infra.doodle.com/v1beta1",
					Kind:       "KafkaTopic",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      kafkaTopicName,
					Namespace: KafkaTopicNamespace,
				},
				Spec: infrav1beta1.KafkaTopicSpec{
					Address:           KafkaBrokersAddress,
					Name:              kafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())

			kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}
			createdKafkaTopic := &infrav1beta1.KafkaTopic{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, createdKafkaTopic)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdKafkaTopic.Spec.Name).Should(Equal(kafkaTopicName))

			By("By checking that partitions are as expected")
			Expect(*createdKafkaTopic.Spec.Partitions).Should(Equal(partitions))

			By("By checking that replication factor is as expected")
			Expect(*createdKafkaTopic.Spec.ReplicationFactor).Should(Equal(replicationFactor))

			By("By checking that kafka brokers address is as expected")
			Expect(createdKafkaTopic.Spec.Address).Should(Equal(KafkaBrokersAddress))

			By("By checking that condition status is true")
			Eventually(func() (metav1.ConditionStatus, error) {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, createdKafkaTopic)
				if err != nil {
					return "", err
				}
				if len(createdKafkaTopic.Status.Conditions) == 0 {
					return "", errors.New("conditions are 0")
				}
				return createdKafkaTopic.Status.Conditions[0].Status, nil
			}, timeout, interval).Should(Equal(metav1.ConditionTrue))

			By("By checking that condition type is Ready")
			Expect(createdKafkaTopic.Status.Conditions[0].Type).Should(Equal(infrav1beta1.ReadyCondition))

			By("By checking that condition reason is Ready")
			Expect(createdKafkaTopic.Status.Conditions[0].Type).Should(Equal(infrav1beta1.ReadyCondition))

			By("By creating a kafka topic in brokers")
			Eventually(func() bool {
				return kafka.DefaultMockKafkaBrokers.GetTopic(kafkaTopicName) != nil
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				topic := kafka.DefaultMockKafkaBrokers.GetTopic(kafkaTopicName)
				if topic == nil {
					return false
				}
				return topic.Name == kafkaTopicName
			}, timeout, interval).Should(BeTrue())

			By("By checking that kafka topic in brokers has expected number of partitions")
			Eventually(func() bool {
				topic := kafka.DefaultMockKafkaBrokers.GetTopic(kafkaTopicName)
				if topic == nil {
					return false
				}
				return topic.Partitions == partitions
			}).Should(BeTrue())

			By("By checking that kafka topic in brokers has expected replication factor")
			Eventually(func() bool {
				topic := kafka.DefaultMockKafkaBrokers.GetTopic(kafkaTopicName)
				if topic == nil {
					return false
				}
				return topic.ReplicationFactor == replicationFactor
			}).Should(BeTrue())
		})
	})

	Context("When decreasing number of partitions", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 3
		kafkaTopicName := "test-decrease-number-of-partitions"
		kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}

		It("Should create new topic", func() {
			By("By checking the topic is created")
			ctx := context.Background()
			kafkaTopic := &infrav1beta1.KafkaTopic{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "kafka.infra.doodle.com/v1beta1",
					Kind:       "KafkaTopic",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      kafkaTopicName,
					Namespace: KafkaTopicNamespace,
				},
				Spec: infrav1beta1.KafkaTopicSpec{
					Address:           KafkaBrokersAddress,
					Name:              kafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())
		})

		It("Should refuse to decrease number of partitions ", func() {
			By("By updating the number of partitions in KafkaTopic object")
			latest := &infrav1beta1.KafkaTopic{}
			var newPartitions int64 = 4
			Eventually(func() error {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, latest)
				if err != nil {
					return err
				}
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				latest.Spec.Partitions = &newPartitions
				return k8sClient.Update(ctx, latest)
			}, timeout, interval).Should(Succeed())

			By("By checking that topic is not ready")
			Eventually(func() error {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, latest)
				if err != nil {
					return err
				}
				if *latest.Spec.Partitions != newPartitions {
					return errors.New("partitions are not changed")
				}
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				if latest.Status.Conditions[0].Status != metav1.ConditionFalse {
					return errors.New("Condition is true")
				}
				return nil
			}, timeout, interval).Should(Succeed())
			By("By checking that reason is that partitions cannot be removed")
			Expect(latest.Status.Conditions[0].Reason).Should(Equal(infrav1beta1.PartitionsFailedToRemoveReason))
		})
	})

	Context("When increasing number of partitions", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 3
		kafkaTopicName := "test-increase-number-of-partitions"
		kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}

		It("Should create new topic", func() {
			By("By checking the topic is created")
			ctx := context.Background()
			kafkaTopic := &infrav1beta1.KafkaTopic{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "kafka.infra.doodle.com/v1beta1",
					Kind:       "KafkaTopic",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      kafkaTopicName,
					Namespace: KafkaTopicNamespace,
				},
				Spec: infrav1beta1.KafkaTopicSpec{
					Address:           KafkaBrokersAddress,
					Name:              kafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())
		})

		It("Should assign new partitions ", func() {
			By("By updating the number of partitions in KafkaTopic object")
			latest := &infrav1beta1.KafkaTopic{}
			var newPartitions int64 = 18
			Eventually(func() error {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, latest)
				if err != nil {
					return err
				}
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				latest.Spec.Partitions = &newPartitions
				return k8sClient.Update(ctx, latest)
			}, timeout, interval).Should(Succeed())

			By("By checking that the number of partitions is properly updated")
			Eventually(func() error {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, latest)
				if err != nil {
					return err
				}
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				if *latest.Spec.Partitions != newPartitions {
					return errors.New("partitions are not changed")
				}
				return nil
			}, timeout, interval).Should(Succeed())
			Eventually(func() int64 {
				updatedTopic := kafka.DefaultMockKafkaBrokers.GetTopic(kafkaTopicName)
				return updatedTopic.Partitions
			}, timeout, interval).Should(Equal(newPartitions))
		})
	})

	Context("When updating topic configuration", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 3
		kafkaTopicName := "test-update-configuration"
		kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}

		It("Should create new topic", func() {
			By("By checking the topic is created")
			ctx := context.Background()
			kafkaTopic := &infrav1beta1.KafkaTopic{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "kafka.infra.doodle.com/v1beta1",
					Kind:       "KafkaTopic",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      kafkaTopicName,
					Namespace: KafkaTopicNamespace,
				},
				Spec: infrav1beta1.KafkaTopicSpec{
					Address:           KafkaBrokersAddress,
					Name:              kafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())
		})

		It("Should update flush.ms", func() {
			By("By checking that flush.ms is properly updated")
			latest := &infrav1beta1.KafkaTopic{}
			var newFlusMs int64 = 666
			Eventually(func() error {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, latest)
				if err != nil {
					return err
				}
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				latest.Spec.KafkaTopicConfig = &infrav1beta1.KafkaTopicConfig{
					FlushMs: &newFlusMs,
				}
				err = k8sClient.Update(ctx, latest)
				if err != nil {
					return err
				}
				return nil
			}, timeout, interval).Should(Succeed())
			Eventually(func() string {
				updatedTopic := kafka.DefaultMockKafkaBrokers.GetTopic(kafkaTopicName)
				return updatedTopic.Config["flush.ms"]
			}, timeout, interval).Should(Equal(strconv.FormatInt(newFlusMs, 10)))
		})
	})

	Context("When updating replication factor for already existing topic", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 3
		kafkaTopicName := "test-update-replication-factor"
		kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}

		It("Should create new topic", func() {
			By("By checking the topic is created")
			ctx := context.Background()
			kafkaTopic := &infrav1beta1.KafkaTopic{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "kafka.infra.doodle.com/v1beta1",
					Kind:       "KafkaTopic",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      kafkaTopicName,
					Namespace: KafkaTopicNamespace,
				},
				Spec: infrav1beta1.KafkaTopicSpec{
					Address:           KafkaBrokersAddress,
					Name:              kafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())
		})

		It("Should refuse to change replication factor", func() {
			By("By updating replication factor on KafkaTopic object")
			latest := &infrav1beta1.KafkaTopic{}
			var newReplicationFactor int64 = 4
			Eventually(func() error {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, latest)
				if err != nil {
					return err
				}
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				latest.Spec.ReplicationFactor = &newReplicationFactor
				return k8sClient.Update(ctx, latest)
			}, timeout, interval).Should(Succeed())

			By("By checking that topic is not ready")
			Eventually(func() error {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, latest)
				if err != nil {
					return err
				}
				if *latest.Spec.ReplicationFactor != newReplicationFactor {
					return errors.New("replication factor is not changed")
				}
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				if latest.Status.Conditions[0].Status != metav1.ConditionFalse {
					return errors.New("Condition is true")
				}
				return nil
			}, timeout, interval).Should(Succeed())
			By("By checking that reason is that replication factor cannot be modified")
			Expect(latest.Status.Conditions[0].Reason).Should(Equal(infrav1beta1.ReplicationFactorFailedToChangeReason))
		})
	})
})
