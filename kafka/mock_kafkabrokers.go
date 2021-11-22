package kafka

import (
	"fmt"
	"github.com/pkg/errors"
)

var DefaultMockKafkaBrokers = MockKafkaBrokers{
	topics: map[string]Topic{},
}

type MockKafkaBrokers struct {
	topics map[string]Topic
}

func (kb *MockKafkaBrokers) GetTopic(name string) *Topic {
	if _, ok := kb.topics[name]; ok {
		topic := kb.topics[name]
		return &topic
	}
	return nil
}

func (kb *MockKafkaBrokers) AddTopic(topic Topic) error {
	if _, ok := kb.topics[topic.Name]; ok {
		return errors.New(fmt.Sprintf("there is already topic named '%s'", topic.Name))
	}
	kb.topics[topic.Name] = topic
	return nil
}

func (kb *MockKafkaBrokers) ClearAllTopics() {
	kb.topics = map[string]Topic{}
}
