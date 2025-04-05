package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type producerImp struct {
	producer *kgo.Client
}

type ProducerConfig struct {
	Brokers string
	AppName string
	Retries int
}

const (
	produceMessage = "produce_message"
	success        = "success"
	failure        = "error"
)

func (p *producerImp) Ping(ctx context.Context) error {
	err := p.producer.Ping(ctx)

	if err != nil {
		errMsg := fmt.Errorf("kafka_producer: broker is not healthy %v", err.Error())

		fmt.Println(errMsg)

		return err
	}

	return nil
}
