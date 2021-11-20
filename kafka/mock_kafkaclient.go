package kafka

import "context"

type MockKafkaClient struct{}

// TODO implement
func NewMockKafkaClient() MockKafkaClient {
	return MockKafkaClient{}
}

func (kc MockKafkaClient) CreateTopic(address string, topic Topic) error {
	return nil
}
func (kc MockKafkaClient) CreatePartitions(ctx context.Context, address string, topic Topic, numberOfPartitions int64) error {
	return nil
}
func (kc MockKafkaClient) UpdateTopicConfiguration(ctx context.Context, address string, topic Topic) error {
	return nil
}
func (kc MockKafkaClient) GetTopic(address string, name string) (*Topic, error) {
	return nil, nil
}
