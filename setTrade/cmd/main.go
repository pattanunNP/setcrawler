package main

import (
	"fmt"
	"log"
	"net/http"
	"settrade/pkg/company"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	langs := []string{"th", "en"}
	var allCompanyData []company.CompanyData

	for _, lang := range langs {
		mainURL := fmt.Sprintf("https://www.set.or.th/%s/market/information/member-list/main", lang)
		companyURLs := fetchCompanyLinks(mainURL)

		if len(companyURLs) == 0 {
			fmt.Printf("No company URLs found for language: %s\n", lang)
			continue
		}

		for _, url := range companyURLs {
			companyData := company.FetchCompanyData(url)
			allCompanyData = append(allCompanyData, companyData)
		}
	}
	company.SaveToJSON("companies_data.json", allCompanyData)
}

func fetchCompanyLinks(url string) []string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Received non-200 response code: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	var companyURLs []string
	doc.Find("div.py-3 div.card.market-related-info a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			fullURL := "https://www.set.or.th" + link
			companyURLs = append(companyURLs, fullURL)
			fmt.Printf("Found URL: %s\n", fullURL)
		}
	})
	return companyURLs
}
