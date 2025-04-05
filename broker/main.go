package broker

import (
	"context"
)

type Broker struct {
	Consumer Consumer
	Producer Producer
}

type Event struct {
	Headers map[string]string
	Values  []byte
	Event   interface{}
}

type Consumer interface {
	// ReadEvent(ctx context.Context) (*Event, error)
	// ReadEventAsync(ctx context.Context, eventChan chan []*Event, errorChan chan error) error
	// CommitMessage(event interface{}) error
	Ping(ctx context.Context) error
}

type Producer interface {
	// SendEvent(ctx context.Context, topic string, channelId string, requestId string, body interface{}) error
	Ping(ctx context.Context) error
}

func NewBroker(consumer Consumer, producer Producer) *Broker {
	return &Broker{
		Consumer: consumer,
		Producer: producer,
	}
}

const (
	CONSUME_MESSAGE = "consume_message"
	PRODUCE_MESSAGE = "produce_message"
	METRIC_ERROR    = "error"
	METRIC_SUCCESS  = "success"
)
