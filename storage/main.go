package main

import (
	"storage/service"
)

func main() {
	service.InitializeRabbitMQ()
	service.InitializeStore()
	service.ActivateRabbitmqConsumer()
}
