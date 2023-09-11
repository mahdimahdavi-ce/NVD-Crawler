package model

import amqp "github.com/rabbitmq/amqp091-go"

type RabbitService struct {
	Channel *amqp.Channel
}
