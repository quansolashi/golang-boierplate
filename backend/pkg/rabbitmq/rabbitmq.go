package rabbitmq

import (
	"context"
	"fmt"
	"os"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Params struct {
	Host     string
	Port     string
	Username string
	Password string
}

type QueueParams struct {
	QueueName string
}

type ExchangeType string

const (
	// ExchangeTypeDirect is used for routing messages to a specific queue.
	ExchangeTypeDirect ExchangeType = "direct"
	// ExchangeTypeFanout is used for broadcasting messages to all queues.
	ExchangeTypeFanout ExchangeType = "fanout"
	// ExchangeTypeTopic is used for routing messages based on wildcard patterns.
	ExchangeTypeTopic ExchangeType = "topic"
	// ExchangeTypeHeaders is used for routing messages based on header attributes.
	ExchangeTypeHeaders ExchangeType = "header"
)

type ExchangeParams struct {
	ExchangeName string
	ExchangeType ExchangeType
}

type Client interface {
	DeclareExchange(params *ExchangeParams) error
	DeclareQueue(params *QueueParams) (amqp091.Queue, error)
	Publish(ctx context.Context, params *PublishParams) error
	Consume(ctx context.Context, params *ConsumeParams, handler Handler) error
	Close()
}

type client struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

type options struct {
	logger zerolog.Logger
}

type Option func(opts *options)

func WithLogger(logger zerolog.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewRabbitMQ(params *Params, opts ...Option) (Client, error) {
	dopts := &options{
		logger: zerolog.New(os.Stdout),
	}
	for i := range opts {
		opts[i](dopts)
	}

	//nolint:nosprintfhostport // do not needed
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", params.Username, params.Password, params.Host, params.Port)
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// ensure worker processed and acknowledged only one message before handling next one
	err = channel.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return nil, err
	}

	return &client{
		conn:    conn,
		channel: channel,
	}, nil
}

func (c *client) DeclareQueue(params *QueueParams) (amqp091.Queue, error) {
	queue, err := c.channel.QueueDeclare(
		params.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return amqp091.Queue{}, err
	}
	return queue, nil
}

func (c *client) DeclareExchange(params *ExchangeParams) error {
	err := c.channel.ExchangeDeclare(
		params.ExchangeName,
		string(params.ExchangeType),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) Close() {
	c.channel.Close()
	c.conn.Close()
}
