package main

import (
	"fmt"
	"scrapper/service"
)

func main() {
	links := service.ExtractVulnerabilitiesLinks("redis")
	details := service.ExtractVulnerabilitiesDetails(links)
	for _, vulner := range details {
		fmt.Println(vulner.CVEID)
		fmt.Println(vulner.Description)
		fmt.Println(vulner.PublishedDate)
		fmt.Println(vulner.LastModified)
		fmt.Println(vulner.NVDScore)
		fmt.Println(vulner.CNAScore)
		for _, val := range vulner.VulnerableVersions {
			fmt.Println(val)
		}
		fmt.Println("***")
	}
}
