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

func ExtractVulnerabilitiesDetails(vulnerabilitiesLinks []string) []Vulnerability {
	fmt.Println("Hit function")
	c := colly.NewCollector()
	vulnerSlice := make([]Vulnerability, len(vulnerabilitiesLinks))
	index := 0

	c.OnHTML("div.col-lg-9:nth-child(1) > p:nth-child(3)", func(h *colly.HTMLElement) {
		// fmt.Println("Hit")
		vulnerSlice[index].description = h.Text
	})

	c.OnHTML("div.col-lg-9:nth-child(1) > p:nth-child(2)", func(h *colly.HTMLElement) {
		// fmt.Println("Hit")
		vulnerSlice[index].description = h.Text
	})

	c.OnHTML("div.bs-callout:nth-child(1)", func(h *colly.HTMLElement) {
		vulnerSlice[index].CVEID = h.ChildText("a")
	})

	c.OnHTML("div.bs-callout:nth-child(1) > span:nth-child(8)", func(h *colly.HTMLElement) {
		vulnerSlice[index].publishedDate = h.Text
	})

	c.OnHTML("div.bs-callout:nth-child(1) > span:nth-child(12)", func(h *colly.HTMLElement) {
		vulnerSlice[index].lastModified = h.Text
	})

	c.OnScraped(func(r *colly.Response) {
		index += 1
		if index < len(vulnerabilitiesLinks) {
			c.Visit(vulnerabilitiesLinks[index])
		}
	})
	fmt.Println(vulnerabilitiesLinks[0])
	c.Visit(vulnerabilitiesLinks[0])

	return vulnerSlice
}
