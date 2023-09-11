package main

import (
	"scrapper/service"
)

func main() {
	query := "mongodb"
	links := service.ExtractVulnerabilitiesLinks(query)
	vulnerabilities := service.ExtractVulnerabilitiesDetails(query, links)
	// for _, vulner := range vulnerabilities {
	// 	fmt.Println(vulner.CVEID)
	// 	fmt.Println(vulner.Description)
	// 	fmt.Println(vulner.PublishedDate)
	// 	fmt.Println(vulner.LastModified)
	// 	fmt.Println(vulner.NVDScore)
	// 	fmt.Println(vulner.CNAScore)
	// 	for _, val := range vulner.VulnerableVersions {
	// 		fmt.Println(val)
	// 	}
	// 	fmt.Println("***")
	// }
	// result := []model.Vulnerability{}
	// result = append(result, model.Vulnerability{
	// 	CVEID:              "5",
	// 	PublishedDate:      "01-01-2012",
	// 	LastModified:       "01-01-2012",
	// 	Description:        "",
	// 	VulnerableVersions: []string{},
	// 	NVDScore:           "",
	// 	CNAScore:           "",
	// })
	service.InitializeRabbitMQ()
	service.SendVulnerabilitiesToDatabase(vulnerabilities)
}
