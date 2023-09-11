package service

import (
	"encoding/json"
	"fmt"
	"log"
	"storage/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitService = new(model.RabbitService)
)

func InitializeRabbitMQ() {
	conn, connErr := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if connErr != nil {
		log.Panicf("Failed to establish a connection to the rabbitMQ server: %s", connErr)
	}

	channel, chanErr := conn.Channel()
	if chanErr != nil {
		log.Panicf("Failed to open a channel to the rabbitMQ server: %s", chanErr)
	}

	rabbitService.Channel = channel
}

func ActivateRabbitmqConsumer() {
	exchErr := rabbitService.Channel.ExchangeDeclare(
		"vulnerability_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil)

	if exchErr != nil {
		log.Panicf("Failed to declare an exchange: %s", exchErr)
	}

	queue, queueErr := rabbitService.Channel.QueueDeclare(
		"vulnerability_queue",
		true,
		false,
		false,
		false,
		nil)

	if queueErr != nil {
		log.Panicf("Failed to create a queue: %s", queueErr)
	}

	bindErr := rabbitService.Channel.QueueBind(
		queue.Name,
		"vulnerability",
		"vulnerability_exchange",
		false,
		nil)

	if bindErr != nil {
		log.Panicf("Failed to bind queue to exchange: %s", bindErr)
	}

	msgs, err := rabbitService.Channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Panicf("Failed to establish consumer: %s", err)
	}

	var forever chan struct{}

	go func() {
		for msg := range msgs {
			vulnerability := new(model.Vulnerability)
			decodeErr := json.Unmarshal(msg.Body, vulnerability)
			if decodeErr != nil {
				fmt.Println(decodeErr)
			}
			acknowledgment := SaveVulnerability(vulnerability)
			if acknowledgment {
				msg.Ack(false)
			}
		}
	}()

	<-forever
}
