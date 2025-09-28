package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn     *amqp091.Connection
	channel  *amqp091.Channel
	exchange string
}



func NewPublisher(amqpURL, exchange string) (*Publisher, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declaramos un exchange tipo "topic" para manejar diferentes eventos
	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &Publisher{conn: conn, channel: ch, exchange: exchange}, nil
}

func (p *Publisher) Publish(routingKey string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return p.channel.Publish(
		p.exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) Close() {
	p.channel.Close()
	p.conn.Close()
}
