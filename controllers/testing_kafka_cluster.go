package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/testcontainers/testcontainers-go"
)

const (
	CLUSTER_NETWORK_NAME = "kafka-cluster"
	ZOOKEEPER_PORT       = "2181"
	KAFKA_BROKER_PORT    = "9092"
	KAFKA_CLIENT_PORT    = "9093"
	KAFKA_BROKER_ID      = "1"
	ZOOKEEPER_IMAGE      = "confluentinc/cp-zookeeper:5.5.2"
	KAFKA_IMAGE          = "confluentinc/cp-kafka:5.5.2"
)

type TestingKafkaCluster struct {
	kafkaContainer     testcontainers.Container
	zookeeperContainer testcontainers.Container
}

func (kc *TestingKafkaCluster) StartCluster() error {
	ctx := context.Background()
	if err := kc.zookeeperContainer.Start(ctx); err != nil {
		return err
	}
	if err := kc.kafkaContainer.Start(ctx); err != nil {
		return err
	}
	if err := kc.createProbe(); err != nil {
		return err
	}
	return kc.startKafka()
}

func (kc *TestingKafkaCluster) StopCluster() error {
	ctx := context.Background()
	if err := kc.kafkaContainer.Terminate(ctx); err != nil {
		return err
	}
	if err := kc.zookeeperContainer.Terminate(ctx); err != nil {
		return err
	}
	return nil
}

func (kc *TestingKafkaCluster) GetKafkaHost() (string, error) {
	ctx := context.Background()
	host, err := kc.kafkaContainer.Host(ctx)
	if err != nil {
		return "", err
	}
	port, err := kc.kafkaContainer.MappedPort(ctx, KAFKA_CLIENT_PORT)
	if err != nil {
		return "", err
	}
	return host + ":" + port.Port(), nil
}

func (kc *TestingKafkaCluster) createProbe() error {
	probeScript, err := ioutil.TempFile("", "probe.sh")
	if err != nil {
		return err
	}
	defer os.Remove(probeScript.Name())
	probeScript.WriteString("#!/bin/bash \n")
	probeScript.WriteString(fmt.Sprintf("id=$(zookeeper-shell localhost:%s ls /brokers/ids | grep \"\\[%s\") \n", ZOOKEEPER_PORT, KAFKA_BROKER_ID))
	probeScript.WriteString(fmt.Sprintf("if [ $id = \"[%s]\" ]; then exit 0; else exit 1; fi", KAFKA_BROKER_ID))

	return kc.zookeeperContainer.CopyFileToContainer(context.Background(), probeScript.Name(), "probe.sh", 0700)
}

func (kc *TestingKafkaCluster) startKafka() error {
	ctx := context.Background()
	kafkaStartFile, err := ioutil.TempFile("", "testcontainers_start.sh")
	if err != nil {
		return err
	}
	defer os.Remove(kafkaStartFile.Name())

	exposedHost, err := kc.GetKafkaHost()
	if err != nil {
		return err
	}
	kafkaStartFile.WriteString("#!/bin/bash \n")
	kafkaStartFile.WriteString(fmt.Sprintf("export KAFKA_ADVERTISED_LISTENERS='PLAINTEXT://%s,BROKER://kafka:%s'\n", exposedHost, KAFKA_BROKER_PORT))
	kafkaStartFile.WriteString(". /etc/confluent/docker/bash-config \n")
	kafkaStartFile.WriteString("/etc/confluent/docker/configure \n")
	kafkaStartFile.WriteString("/etc/confluent/docker/launch \n")

	return kc.kafkaContainer.CopyFileToContainer(ctx, kafkaStartFile.Name(), "testcontainers_start.sh", 0700)
}

func NewTestingKafkaCluster() (*TestingKafkaCluster, error) {
	ctx := context.Background()

	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Name: CLUSTER_NETWORK_NAME},
	})
	if err != nil {
		return nil, err
	}
	dockerNetwork := network.(*testcontainers.DockerNetwork)
	zookeeperContainer, err := createZookeeperContainer(dockerNetwork)
	if err != nil {
		return nil, err
	}
	kafkaContainer, err := createKafkaContainer(dockerNetwork)
	if err != nil {
		return nil, err
	}
	return &TestingKafkaCluster{
		kafkaContainer:     kafkaContainer,
		zookeeperContainer: zookeeperContainer,
	}, nil
}

func (kc *TestingKafkaCluster) IsAlive() (bool, error) {
	if s, _, err := kc.zookeeperContainer.Exec(context.Background(), []string{"/probe.sh"}); err != nil {
		return false, err
	} else if s == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func createZookeeperContainer(network *testcontainers.DockerNetwork) (testcontainers.Container, error) {
	ctx := context.Background()

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:          ZOOKEEPER_IMAGE,
			ExposedPorts:   []string{ZOOKEEPER_PORT},
			Env:            map[string]string{"ZOOKEEPER_CLIENT_PORT": ZOOKEEPER_PORT, "ZOOKEEPER_TICK_TIME": "2000"},
			Networks:       []string{network.Name},
			NetworkAliases: map[string][]string{network.Name: {"zookeeper"}},
		},
	})
}

func createKafkaContainer(network *testcontainers.DockerNetwork) (testcontainers.Container, error) {
	ctx := context.Background()

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        KAFKA_IMAGE,
			ExposedPorts: []string{KAFKA_CLIENT_PORT},
			Env: map[string]string{
				"KAFKA_BROKER_ID":                        KAFKA_BROKER_ID,
				"KAFKA_ZOOKEEPER_CONNECT":                fmt.Sprintf("zookeeper:%s", ZOOKEEPER_PORT),
				"KAFKA_LISTENERS":                        fmt.Sprintf("PLAINTEXT://0.0.0.0:%s,BROKER://0.0.0.0:%s", KAFKA_CLIENT_PORT, KAFKA_BROKER_PORT),
				"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":   "BROKER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
				"KAFKA_INTER_BROKER_LISTENER_NAME":       "BROKER",
				"KAFKA_OFFESTS_TOPIC_REPLICATION_FACTOR": "1",
			},
			Networks:       []string{network.Name},
			NetworkAliases: map[string][]string{network.Name: {"kafka"}},
			Cmd:            []string{"sh", "-c", "while [ ! -f /testcontainers_start.sh ]; do sleep 0.1; done; /testcontainers_start.sh"},
		},
	})
}
