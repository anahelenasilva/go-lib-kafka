package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
)

type kafkaConsumerImp struct {
	consumer    *kgo.Client
	topic       string
	poolrecords int
}

type KafkaConsumerConfig struct {
	Broker      string
	GroupId     string
	AutoOffset  string
	Topic       string
	AutoCommit  *bool
	Poolrecords int
}

const (
	consumeMessage = "consume_message"
)

func (c *KafkaConsumerConfig) validate() error {
	if c.Broker == "" {
		return errors.New("kafka_consumer: broker is required")
	}
	if c.GroupId == "" {
		return errors.New("kafka_consumer: groupId is required")
	}
	if c.Topic == "" {
		return errors.New("kafka_consumer: topic is required")
	}

	return nil
}

func (c *KafkaConsumerConfig) normalize() {
	if c.AutoCommit == nil {
		autoCommit := true
		c.AutoCommit = &autoCommit
	}

	if c.GroupId == "" {
		c.GroupId = "default_consumer_group"
	}

	if c.Poolrecords == 0 {
		c.Poolrecords = 100
	}
}

func NewKafkaConsumer(config *KafkaConsumerConfig) (*kafkaConsumerImp, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}

	config.normalize()

	opts := []kgo.Opt{
		kgo.SeedBrokers(strings.Split(config.Broker, ",")...),
		kgo.ConsumerGroup(config.GroupId),
		kgo.ConsumeTopics(config.Topic),
	}

	if config.AutoOffset == "latest" {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()))
	} else if config.AutoOffset == "earliest" {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()))
	}
	if !*config.AutoCommit {
		opts = append(opts, kgo.DisableAutoCommit())
	}

	kgoConsumerClient, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumerImp{
		consumer:    kgoConsumerClient,
		topic:       config.Topic,
		poolrecords: config.Poolrecords,
	}, nil
}

func (p *kafkaConsumerImp) Ping(ctx context.Context) error {
	err := p.consumer.Ping(ctx)

	if err != nil {
		errMsg := fmt.Errorf("kafka_consumer: failed to ping kafka. error: %s", err.Error())

		fmt.Println(errMsg)

		return err
	}

	return nil
}
