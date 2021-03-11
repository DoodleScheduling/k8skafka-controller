package kafka

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	k "github.com/segmentio/kafka-go"
	"math"
	"net"
	"strconv"
	"time"
)

const (
	TCP = "tcp"
)

const (
	DefaultKafkaClientTimeout = 4 * time.Minute
)

type Connection struct {
	protocol string
	uri      string
}

type Topic struct {
	Name              string
	Partitions        int64
	ReplicationFactor int64
	Brokers           map[int64]int64
}

func NewTCPConnection(uri string) *Connection {
	return NewConnection(TCP, uri)
}

func NewConnection(protocol string, uri string) *Connection {
	return &Connection{
		protocol: protocol,
		uri:      uri,
	}
}

func (c *Connection) CreateTopic(topic Topic) error {
	topicConfig := k.TopicConfig{
		Topic:             topic.Name,
		NumPartitions:     int(topic.Partitions),
		ReplicationFactor: int(topic.ReplicationFactor),
	}

	conn, err := k.Dial(c.protocol, c.uri)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot dial %s via %s", c.uri, c.protocol))
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot get controller for %s", c.uri))
	}
	var controllerConn *k.Conn
	controllerConn, err = k.Dial(c.protocol, net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot dial controller connection %s %s", controller.Host, strconv.Itoa(controller.Port)))
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot create topic %+v for %s", topicConfig, c.uri))
	}
	return nil
}

func (c *Connection) CreatePartitions(ctx context.Context, topic Topic, numberOfPartitions int64) error {
	addr, err := net.ResolveTCPAddr(c.protocol, c.uri)
	if err != nil {
		return err
	}
	client := k.Client{
		Addr:    addr,
		Timeout: DefaultKafkaClientTimeout,
	}

	topicPartitionAssignments := make([]k.TopicPartitionAssignment, 0)

	var p int64
	brokers := make(map[int64]int64)
	for b := range topic.Brokers {
		brokers[b] = topic.Brokers[b]
	}

	for p = 0; p < numberOfPartitions; p++ {
		var brokerIDs []int32
		brokerIDs, brokers = c.assignBrokersToPartition(topic, brokers)
		// TODO check how partition leader is assigned. Is it the first broker ID in slice? Not having a balanced leader assignment leads to leader skewed scenario
		topicPartitionAssignments = append(topicPartitionAssignments, k.TopicPartitionAssignment{BrokerIDs: brokerIDs})
	}

	topicPartitionsConfig := []k.TopicPartitionsConfig{
		{
			Name:                      topic.Name,
			Count:                     int32(topic.Partitions),
			TopicPartitionAssignments: topicPartitionAssignments,
		},
	}

	req := k.CreatePartitionsRequest{
		Topics:       topicPartitionsConfig,
		ValidateOnly: false,
	}
	res, err := client.CreatePartitions(ctx, &req)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot create partitions via client %s %s", c.uri, c.protocol))
	}
	if res.Errors == nil {
		return nil
	}
	if e, found := res.Errors[topic.Name]; found {
		return errors.Wrap(e, fmt.Sprintf("found error while creating partitions via client %s %s", c.uri, c.protocol))
	}
	return nil
}

func (c *Connection) GetTopic(name string) (*Topic, error) {
	conn, err := k.Dial(c.protocol, c.uri)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("cannot get topic - dial failed for %s via %s", c.uri, c.protocol))
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	var numberOfPartitions int64 = 0
	var numberOfReplicas int64 = 0
	brokers := make(map[int64]int64)

	for _, p := range partitions {
		if p.Topic != name {
			continue
		}
		numberOfPartitions++
		if _, ok := brokers[int64(p.Leader.ID)]; !ok {
			brokers[int64(p.Leader.ID)] = 0
		} else {
			brokers[int64(p.Leader.ID)] = brokers[int64(p.Leader.ID)] + 1
		}
		numberOfReplicas = int64(len(p.Replicas))
	}

	if numberOfPartitions == 0 {
		return nil, nil
	}
	return &Topic{
		Name:              name,
		Partitions:        numberOfPartitions,
		Brokers:           brokers,
		ReplicationFactor: numberOfReplicas,
	}, nil
}

// Assign broker IDs by least used brokers currently
// Return selected brokerIDs, and updated collection of all brokers, with new data about usage
func (c *Connection) assignBrokersToPartition(topic Topic, brokers map[int64]int64) ([]int32, map[int64]int64) {
	brokerIDs := make([]int32, 0)
	var i int64
	for i = 0; i < topic.ReplicationFactor; i++ {
		var selectedBroker int64 = 0
		var min int64 = math.MaxInt64
		for b, v := range brokers {
			alreadySelected := false
			for _, bid := range brokerIDs {
				if int64(bid) == b {
					alreadySelected = true
				}
			}
			if !alreadySelected {
				if v < min {
					min = v
					selectedBroker = b
				}
			}
		}
		brokerIDs = append(brokerIDs, int32(selectedBroker))
		brokers[selectedBroker] = brokers[selectedBroker] + 1
	}
	return brokerIDs, brokers
}
