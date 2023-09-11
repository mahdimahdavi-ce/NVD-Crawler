package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"scrapper/model"

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

func SendVulnerabilitiesToDatabase(vulnerabilities []model.Vulnerability) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, vulnerability := range vulnerabilities {
		jsonData, err := json.Marshal(vulnerability)
		if err != nil {
			fmt.Printf("Failed to marshal json: %v", jsonData)
		}
		rabbitService.Channel.PublishWithContext(
			ctx,
			"vulnerability_exchange",
			"vulnerability",
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/json",
				Body:         []byte(jsonData),
				DeliveryMode: amqp.Persistent,
			},
		)
		fmt.Printf("Vulnerability with the id of %s is send to Storage service\n", vulnerability.CVEID)
	}
}
