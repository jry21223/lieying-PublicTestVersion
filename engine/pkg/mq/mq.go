package mq

import (
	"fmt"

	"github.com/lieying/engine/internal/config"
	"github.com/streadway/amqp"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(cfg config.RabbitMQConfig) (*Client, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%d%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{conn: conn, channel: ch}, nil
}

func (c *Client) DeclareQueue(name string) (amqp.Queue, error) {
	return c.channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *Client) Publish(queue string, body []byte) error {
	return c.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *Client) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *Client) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}
