package kafka

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

type MockKafkaClient struct{}

func NewMockKafkaClient() MockKafkaClient {
	return MockKafkaClient{}
}

func (kc MockKafkaClient) CreateTopic(address string, topic Topic) error {
	return DefaultMockKafkaBrokers.AddTopic(topic)
}
func (kc MockKafkaClient) CreatePartitions(ctx context.Context, address string, topic Topic, numberOfPartitions int64) error {
	topicFromBrokers := DefaultMockKafkaBrokers.GetTopic(topic.Name)
	if topicFromBrokers == nil {
		return errors.New(fmt.Sprintf("topic %s does not exist in the mock brokers", topic.Name))
	}
	if topic.Partitions-topicFromBrokers.Partitions != numberOfPartitions {
		return errors.New(fmt.Sprintf("mismatch between new number of partitions (%d) and supplied partitions to create (%d)", topic.Partitions, numberOfPartitions))
	}
	topicFromBrokers.Partitions = topic.Partitions
	DefaultMockKafkaBrokers.UpdateTopic(*topicFromBrokers)
	return nil
}
func (kc MockKafkaClient) UpdateTopicConfiguration(ctx context.Context, address string, topic Topic) error {
	topicFromBrokers := DefaultMockKafkaBrokers.GetTopic(topic.Name)
	if topicFromBrokers == nil {
		return errors.New(fmt.Sprintf("topic %s does not exist in the mock brokers", topic.Name))
	}
	topicFromBrokers.Config = topic.Config
	DefaultMockKafkaBrokers.UpdateTopic(*topicFromBrokers)
	return nil
}
func (kc MockKafkaClient) GetTopic(address string, name string) (*Topic, error) {
	return DefaultMockKafkaBrokers.GetTopic(name), nil
}
