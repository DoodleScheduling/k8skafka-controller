package kafka

import (
	"context"
	k "github.com/segmentio/kafka-go"
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
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *k.Conn
	controllerConn, err = k.Dial(c.protocol, net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) CreatePartitions(ctx context.Context, topic Topic) error {
	addr, err := net.ResolveTCPAddr(c.protocol, c.uri)
	if err != nil {
		return err
	}
	client := k.Client{
		Addr:    addr,
		Timeout: DefaultKafkaClientTimeout,
	}
	topicPartitionsConfig := []k.TopicPartitionsConfig{
		{
			Name:  topic.Name,
			Count: int32(topic.Partitions),
		},
	}
	req := k.CreatePartitionsRequest{
		Topics:       topicPartitionsConfig,
		ValidateOnly: false,
	}
	res, err := client.CreatePartitions(ctx, &req)
	if err != nil {
		return err
	}
	if res.Errors == nil {
		return nil
	}
	if e, found := res.Errors[topic.Name]; found {
		return e
	}
	return nil
}
