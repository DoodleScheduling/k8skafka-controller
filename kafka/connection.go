package kafka

import (
	k "github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

const (
	TCP = "tcp"
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

func (r *Connection) CreateTopic(topic Topic) error {
	topicConfig := k.TopicConfig{
		Topic:             topic.Name,
		NumPartitions:     int(topic.Partitions),
		ReplicationFactor: int(topic.ReplicationFactor),
	}

	conn, err := k.Dial(r.protocol, r.uri)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *k.Conn
	controllerConn, err = k.Dial(r.protocol, net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
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
