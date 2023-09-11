package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"scrapper/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

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
		fmt.Printf("New vulnerability link is found: %v\n", vulnerability)
	})

	c.Visit(link)

	return vulnerabilitiesLinks
}

func ExtractVulnerabilitiesDetails(query string, vulnerabilitiesLinks []string) []model.Vulnerability {
	c := colly.NewCollector()

	vulnerSlice := make([]model.Vulnerability, len(vulnerabilitiesLinks))
	index := 0
	// Extract the description of vulnerability
	c.OnHTML("div.col-lg-9:nth-child(1) > p:nth-child(3)", func(h *colly.HTMLElement) {
		vulnerSlice[index].Description = h.Text
	})
	// Extract the description of vulnerability if the first one dosen't work
	c.OnHTML("div.col-lg-9:nth-child(1) > p:nth-child(2)", func(h *colly.HTMLElement) {
		vulnerSlice[index].Description = h.Text
	})
	// Extract the CVE_ID of vulnerabilty
	c.OnHTML("div.bs-callout:nth-child(1)", func(h *colly.HTMLElement) {
		vulnerSlice[index].CVEID = h.ChildText("a")
	})
	// Extract the publish date of vulnerabilty
	c.OnHTML("div.bs-callout:nth-child(1) > span:nth-child(8)", func(h *colly.HTMLElement) {
		vulnerSlice[index].PublishedDate = h.Text
	})
	// Extract the last modified date of vulnerabilty
	c.OnHTML("div.bs-callout:nth-child(1) > span:nth-child(12)", func(h *colly.HTMLElement) {
		vulnerSlice[index].LastModified = h.Text
	})
	// Extract the NVD severity score of vulnerabilty
	c.OnHTML("#Cvss3NistCalculatorAnchor", func(h *colly.HTMLElement) {
		vulnerSlice[index].NVDScore = h.Text
	})
	// Extract the CNA severity score of vulnerabilty
	c.OnHTML("#Cvss3CnaCalculatorAnchor", func(h *colly.HTMLElement) {
		vulnerSlice[index].CNAScore = h.Text
	})
	// Extract vulnerable versions
	c.OnHTML("body", func(h *colly.HTMLElement) {
		res, err := http.Get(vulnerabilitiesLinks[index])
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			fmt.Println("Error creating GoQuery document:", err)
			return
		}

		result := []string{}
		doc.Find("td[data-testid*='vuln-change-history']").Each(func(i int, e *goquery.Selection) {
			// if it's a td tag that inculdes a portion of the vulnerable versions
			if strings.Contains(e.Text(), "*cpe") {
				extractVulnerableVersions(e.Text(), &result)
			}
		})
		vulnerSlice[index].Name = query
		vulnerSlice[index].VulnerableVersions = result
	})

	c.OnScraped(func(r *colly.Response) {
		index += 1
		if index < len(vulnerabilitiesLinks) {
			c.Visit(vulnerabilitiesLinks[index])
		}
	})

	c.Visit(vulnerabilitiesLinks[0])

	return vulnerSlice
}

func splitBeforeSeparator(input, separator string) []string {
	var result []string
	parts := strings.Split(input, separator)

	for i, part := range parts {
		if i > 0 {
			part = separator + part
		}
		result = append(result, part)
	}

	return result
}

func extractVulnerableVersions(elementText string, result *[]string) {
	strSlice := splitBeforeSeparator(elementText, "*cpe")
	for index, str := range strSlice {
		if strings.HasPrefix(str, "*cpe") {
			if strings.Contains(str, "versions") {
				arr := splitBeforeSeparator(str, "versions")
				str = strings.TrimSpace(arr[1])
				if index != len(strSlice)-1 {
					*result = append(*result, str)
				} else {
					s := strings.Split(str, "\n")
					*result = append(*result, string(s[0]))
				}
			}
		}
	}
}
