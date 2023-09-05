package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Vulnerability struct {
	CVEID         string
	publishedDate string
	lastModified  string
	description   string
	severity      string
}

func main() {
	links := ExtractVulnerabilitiesLinks("redis")
	details := ExtractVulnerabilitiesDetails(links)
	for _, vulner := range details {
		fmt.Println(vulner.CVEID)
		fmt.Println(vulner.description)
		fmt.Println(vulner.publishedDate)
		fmt.Println(vulner.lastModified)
		fmt.Println("***")
	}
}

func generateLink(query string) string {
	baseUrl := "https://nvd.nist.gov/vuln/search/results"
	formType := "Basic"
	resultsType := "overview"
	queryType := "phrase"
	searchType := "last3months" // or it could be "all"
	isCpeNameSearch := false
	return fmt.Sprintf("%s?form_type=%v&results_type=%v&query=%v&queryType=%v&search_type=%v&isCpeNameSearch=%v",
		baseUrl,
		formType,
		resultsType,
		query,
		queryType,
		searchType,
		isCpeNameSearch)
}

func ExtractVulnerabilitiesLinks(query string) []string {
	c := colly.NewCollector()
	vulnerabilitiesLinks := []string{}
	link := generateLink(query)

	c.OnHTML("tr th strong", func(h *colly.HTMLElement) {
		vulnerability := h.ChildText("a")
		vulnerabilitiesLinks = append(vulnerabilitiesLinks, fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%v", vulnerability))
		fmt.Printf("New link is found: %v\n", vulnerability)
	})

	c.Visit(link)

	return vulnerabilitiesLinks
}
