package kafka

import "context"

type Topic struct {
	Name              string
	Partitions        int64
	ReplicationFactor int64
	Brokers           Brokers
	Config            Config
}

// Holds BrokerIds as keys, and number of assignments for this Broker (currently the partitions that the Broker is a leader for)
type Brokers map[int64]int64

// Holds name-value configuration options, passed to underlying kafka libs as-is
type Config map[string]string

type KafkaClient interface {
	CreateTopic(uri string, topic Topic) error
	CreatePartitions(ctx context.Context, uri string, topic Topic, numberOfPartitions int64) error
	UpdateTopicConfiguration(ctx context.Context, uri string, topic Topic) error
	GetTopic(uri string, name string) (*Topic, error)
}
