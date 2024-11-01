package rabbitmq

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type PublishParams struct {
	QueueName string
	Message   string
}

func (c *client) Publish(ctx context.Context, params *PublishParams) error {
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

	err = c.channel.PublishWithContext(ctx,
		"",
		params.QueueName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(params.Message),
		},
	)
	if err != nil {
		return err
	}
	fmt.Println("published message to queue:", params.QueueName)

	return nil
}
