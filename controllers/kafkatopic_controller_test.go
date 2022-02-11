package controllers

import (
	"context"
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"strconv"
	"time"

	infrav1beta1 "github.com/DoodleScheduling/k8skafka-controller/api/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("KafkaTopic controller", func() {
	const (
		KafkaTopicNamespace = "default"
		timeout             = time.Second * 20
		interval            = time.Millisecond * 250
	)

	Context("When creating a topic that doesn't exist already", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 1
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
					Address:           TestingKafkaClusterHost,
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
			Expect(createdKafkaTopic.Spec.Address).Should(Equal(TestingKafkaClusterHost))

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

			By("By checking that the topic is created in kafka cluster")
			Eventually(func() error {
				_, err := GetTopic(kafkaTopicName)
				return err
			}, timeout, interval).Should(Succeed())

			topicInCluster, err := GetTopic(kafkaTopicName)
			Expect(err).To(BeNil())
			Expect(topicInCluster).ToNot(BeNil())

			By("By checking that topic in cluster has expected number of partitions")
			Expect(int64(len(topicInCluster.Partitions))).To(Equal(partitions))
			By("By checking that kafka topic in brokers has expected replication factor")
			Eventually(func() bool {
				for _, p := range topicInCluster.Partitions {
					if int64(len(p.Replicas)) != replicationFactor {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})
	})

	Context("When decreasing number of partitions", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 1
		kafkaTopicName := "test-decrease-number-of-partitions"
		kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}

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
					Address:           TestingKafkaClusterHost,
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

			By("By checking that the topic is created in kafka cluster")
			Eventually(func() error {
				_, err := GetTopic(kafkaTopicName)
				return err
			}, timeout, interval).Should(Succeed())

			topicInCluster, err := GetTopic(kafkaTopicName)
			Expect(err).To(BeNil())
			Expect(topicInCluster).ToNot(BeNil())

			By("By checking that topic in cluster has expected number of partitions")
			Expect(int64(len(topicInCluster.Partitions))).To(Equal(partitions))
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
				// Had to remove this one, as it is too flaky in CI environment
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				if latest.Status.Conditions[0].Status != metav1.ConditionFalse {
					return errors.New("condition is true")
				}
				return nil
			}, timeout, interval).Should(Succeed())

			// Had to remove this one, as it is too flaky in CI environment
			By("By checking that reason is that partitions cannot be removed")
			Expect(latest.Status.Conditions[0].Reason).Should(Equal(infrav1beta1.PartitionsFailedToRemoveReason))

			By("By checking that the number of partitions for topic in kafka brokers hasn't changed")
			Eventually(func() error {
				_, err := GetTopic(kafkaTopicName)
				return err
			}, timeout, interval).Should(Succeed())
			topicInCluster, err := GetTopic(kafkaTopicName)
			Expect(err).To(BeNil())
			Expect(topicInCluster).ToNot(BeNil())
			Expect(int64(len(topicInCluster.Partitions))).To(Equal(partitions))
		})
	})

	Context("When increasing number of partitions", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 1
		kafkaTopicName := "test-increase-number-of-partitions"
		kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}

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
					Address:           TestingKafkaClusterHost,
					Name:              kafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())

			createdKafkaTopic := &infrav1beta1.KafkaTopic{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, createdKafkaTopic)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("By checking that the topic is created in kafka cluster")
			Eventually(func() error {
				_, err := GetTopic(kafkaTopicName)
				return err
			}, timeout, interval).Should(Succeed())

			topicInCluster, err := GetTopic(kafkaTopicName)
			Expect(err).To(BeNil())
			Expect(topicInCluster).ToNot(BeNil())

			By("By checking that topic in cluster has expected number of partitions")
			Expect(int64(len(topicInCluster.Partitions))).To(Equal(partitions))
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
			Eventually(func() bool {
				topic, err := GetTopic(kafkaTopicName)
				if err != nil {
					return false
				}
				return int64(len(topic.Partitions)) == newPartitions
			}, timeout, interval).Should(BeTrue())
		})
	})

	Context("When changing replication factor for already existing topic", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 1
		kafkaTopicName := "test-change-replication-factor"
		kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}

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
					Address:           TestingKafkaClusterHost,
					Name:              kafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())

			createdKafkaTopic := &infrav1beta1.KafkaTopic{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, createdKafkaTopic)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("By checking that the topic is created in kafka cluster")
			Eventually(func() error {
				_, err := GetTopic(kafkaTopicName)
				return err
			}, timeout, interval).Should(Succeed())

			topicInCluster, err := GetTopic(kafkaTopicName)
			Expect(err).To(BeNil())
			Expect(topicInCluster).ToNot(BeNil())

			By("By checking that topic in cluster has expected number of partitions")
			Expect(int64(len(topicInCluster.Partitions))).To(Equal(partitions))
		})

		It("Should refuse to change replication factor", func() {
			By("By updating the replication factor in KafkaTopic object")
			latest := &infrav1beta1.KafkaTopic{}
			var newReplicationFactor int64 = 2
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
				// Had to remove this one, as it is too flaky in CI environment
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				if latest.Status.Conditions[0].Status != metav1.ConditionFalse {
					return errors.New("condition is true")
				}
				return nil
			}, timeout, interval).Should(Succeed())

			// Had to remove this one, as it is too flaky in CI environment
			By("By checking that reason is that replication factor cannot be changed")
			Expect(latest.Status.Conditions[0].Reason).Should(Equal(infrav1beta1.ReplicationFactorFailedToChangeReason))

			By("By checking that topic replication factor in kafka brokers hasn't changed")
			Eventually(func() error {
				_, err := GetTopic(kafkaTopicName)
				return err
			}, timeout, interval).Should(Succeed())
			topicInCluster, err := GetTopic(kafkaTopicName)
			Expect(err).To(BeNil())
			Expect(topicInCluster).ToNot(BeNil())
			Eventually(func() bool {
				for _, p := range topicInCluster.Partitions {
					if int64(len(p.Replicas)) != replicationFactor {
						return false
					}
				}
				return true
			}).Should(BeTrue())
		})
	})

	// These series of tests use data extracted to separate struct (KafkaTopicConfigTestData). For each data entry, one test is created
	Context("When updating topic configuration", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 1
		kafkaTopicName := "test-update-configuration"
		kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}

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
					Address:           TestingKafkaClusterHost,
					Name:              kafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())
			createdKafkaTopic := &infrav1beta1.KafkaTopic{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, createdKafkaTopic)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("By checking that the topic is created in kafka cluster")
			Eventually(func() error {
				_, err := GetTopic(kafkaTopicName)
				return err
			}, timeout, interval).Should(Succeed())

			topicInCluster, err := GetTopic(kafkaTopicName)
			Expect(err).To(BeNil())
			Expect(topicInCluster).ToNot(BeNil())
		})

		// Create a test for each entry in KafkaTopicConfigTestData (kafkatopic_controller_test_data.go)
		for configName, testConfigs := range KafkaTopicConfigTestData {
			for _, testConfig := range testConfigs {
				// store into variable in this closure, otherwise next loop iteration will override variable
				cn := configName
				tc := testConfig

				It(fmt.Sprintf("Should update %s", configName), func() {
					By(fmt.Sprintf("By checking that %s is properly updated", configName))
					latest := &infrav1beta1.KafkaTopic{}
					Eventually(func() error {
						err := k8sClient.Get(ctx, kafkaTopicLookupKey, latest)
						if err != nil {
							return err
						}
						if len(latest.Status.Conditions) == 0 {
							return errors.New("conditions are 0")
						}
						latest.Spec.KafkaTopicConfig = tc.createKafkaTopicObjectF(tc.valueToSet)
						err = k8sClient.Update(ctx, latest)
						if err != nil {
							return err
						}
						return nil
					}, timeout, interval).Should(Succeed())
					Eventually(func() string {
						latestConfig, err := GetTopicConfig(kafkaTopicName)
						if err != nil {
							return err.Error()
						}
						for _, c := range latestConfig {
							for _, cc := range c.ConfigEntries {
								if cc.ConfigName == cn {
									return cc.ConfigValue
								}
							}
						}
						return ""
					}, timeout, interval).Should(Equal(tc.expectedValue))
				})
			}
		}
	})

	Context("When updating topic configuration for topic that already has configuration", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 1
		var retentionMs int64 = 1000
		var flushMs int64 = 666
		kafkaTopicName := "test-update-configuration-for-existing-configuration"
		kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}

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
					Address:           TestingKafkaClusterHost,
					Name:              kafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
					KafkaTopicConfig: &infrav1beta1.KafkaTopicConfig{
						RetentionMs: &retentionMs,
					},
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())
			createdKafkaTopic := &infrav1beta1.KafkaTopic{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, createdKafkaTopic)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("By checking that the topic is created in kafka cluster")
			Eventually(func() error {
				_, err := GetTopic(kafkaTopicName)
				return err
			}, timeout, interval).Should(Succeed())

			topicInCluster, err := GetTopic(kafkaTopicName)
			Expect(err).To(BeNil())
			Expect(topicInCluster).ToNot(BeNil())
		})

		It("Should keep existing configuration", func() {
			By("By adding new configuration option")
			latest := &infrav1beta1.KafkaTopic{}
			Eventually(func() error {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, latest)
				if err != nil {
					return err
				}
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				if latest.Spec.KafkaTopicConfig == nil {
					return errors.New("config is nil")
				}
				latest.Spec.KafkaTopicConfig.FlushMs = &flushMs
				err = k8sClient.Update(ctx, latest)
				if err != nil {
					return err
				}
				return nil
			}, timeout, interval).Should(Succeed())

			By("By checking that new configuration option is added to topic in cluster")
			Eventually(func() string {
				existingConfig, err := GetTopicConfig(kafkaTopicName)
				if err != nil {
					return err.Error()
				}
				for _, c := range existingConfig {
					for _, cc := range c.ConfigEntries {
						if cc.ConfigName == FlushMs {
							return cc.ConfigValue
						}
					}
				}
				return ""
			}, timeout, interval).Should(Equal(fmt.Sprint(flushMs)))

			By("By checking that preexisting topic configuration option is still in brokers")
			Eventually(func() string {
				existingConfig, err := GetTopicConfig(kafkaTopicName)
				if err != nil {
					return err.Error()
				}
				for _, c := range existingConfig {
					for _, cc := range c.ConfigEntries {
						if cc.ConfigName == RetentionMs {
							return cc.ConfigValue
						}
					}
				}
				return ""
			}, timeout, interval).Should(Equal(fmt.Sprint(retentionMs)))
		})
	})

	Context("When removing one configuration option from topic", func() {
		var partitions int64 = 16
		var replicationFactor int64 = 1
		var retentionMs int64 = 1000
		var flushMs int64 = 666
		kafkaTopicName := "test-update-configuration-for-removing-configuration"
		kafkaTopicLookupKey := types.NamespacedName{Name: kafkaTopicName, Namespace: KafkaTopicNamespace}

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
					Address:           TestingKafkaClusterHost,
					Name:              kafkaTopicName,
					Partitions:        &partitions,
					ReplicationFactor: &replicationFactor,
					KafkaTopicConfig: &infrav1beta1.KafkaTopicConfig{
						RetentionMs: &retentionMs,
						FlushMs:     &flushMs,
					},
				},
			}
			Expect(k8sClient.Create(ctx, kafkaTopic)).Should(Succeed())
			createdKafkaTopic := &infrav1beta1.KafkaTopic{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, createdKafkaTopic)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("By checking that the topic is created in kafka cluster")
			Eventually(func() error {
				_, err := GetTopic(kafkaTopicName)
				return err
			}, timeout, interval).Should(Succeed())

			topicInCluster, err := GetTopic(kafkaTopicName)
			Expect(err).To(BeNil())
			Expect(topicInCluster).ToNot(BeNil())
		})

		It("Should remove one configuration option", func() {
			By("By removing a configuration option")
			latest := &infrav1beta1.KafkaTopic{}
			Eventually(func() error {
				err := k8sClient.Get(ctx, kafkaTopicLookupKey, latest)
				if err != nil {
					return err
				}
				if len(latest.Status.Conditions) == 0 {
					return errors.New("conditions are 0")
				}
				if latest.Spec.KafkaTopicConfig == nil {
					return errors.New("config is nil")
				}
				latest.Spec.KafkaTopicConfig = &infrav1beta1.KafkaTopicConfig{
					RetentionMs: &retentionMs,
				}
				err = k8sClient.Update(ctx, latest)
				if err != nil {
					return err
				}
				return nil
			}, timeout, interval).Should(Succeed())

			By("By checking that removed configuration option is removed from topic in cluster")
			Eventually(func() bool {
				latestConfig, err := GetTopicConfig(kafkaTopicName)
				if err != nil {
					return false
				}
				for _, c := range latestConfig {
					for _, cc := range c.ConfigEntries {
						if cc.ConfigName == FlushMs {
							if cc.ConfigValue == strconv.FormatInt(flushMs, 10) {
								return false
							}
						}
					}
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("By checking that other topic configuration option is still in brokers")
			Eventually(func() string {
				latestConfig, err := GetTopicConfig(kafkaTopicName)
				if err != nil {
					return err.Error()
				}
				for _, c := range latestConfig {
					for _, cc := range c.ConfigEntries {
						if cc.ConfigName == RetentionMs {
							return cc.ConfigValue
						}
					}
				}
				return ""
			}, timeout, interval).Should(Equal(fmt.Sprint(retentionMs)))
		})
	})
})
