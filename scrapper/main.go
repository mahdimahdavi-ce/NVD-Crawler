package main

import (
	"scrapper/service"
)

func main() {
	query := "mongodb"
	links := service.ExtractVulnerabilitiesLinks(query)
	vulnerabilities := service.ExtractVulnerabilitiesDetails(query, links)
	service.InitializeRabbitMQ()
	service.SendVulnerabilitiesToDatabase(vulnerabilities)
}
