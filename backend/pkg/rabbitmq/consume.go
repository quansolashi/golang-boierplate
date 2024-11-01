package rabbitmq

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type ConsumeParams struct {
	QueueName string
}

type Handler func(message amqp091.Delivery)

func (c *client) Consume(ctx context.Context, params *ConsumeParams, handler Handler) error {
	_, err := c.channel.QueueDeclare(
		params.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	fmt.Println("consume with queue:", params.QueueName)

	messages, err := c.channel.ConsumeWithContext(ctx,
		params.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			handler(d) // handle message
		}
	}()
	<-forever

	return nil
}
