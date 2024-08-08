package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"setResearch/client"
	"setResearch/parser"
)

func main() {
	baseURL := "https://www.settrade.com/api/cms/v1/research-settrade/search?startDate=01%2F08%2F2024&endDate=08%2F08%2F2024&pageSize=10"
	maxPages := 3
	allResponses := []client.Response{}

	for pageIndex := 0; pageIndex < maxPages; pageIndex++ {
		response, err := client.FetchPage(baseURL, pageIndex)
		if err != nil {
			log.Fatalf("error fetching page %d: %v", pageIndex, err)
		}

		for i := range response.ResearchItems.Items {
			item := &response.ResearchItems.Items[i]
			htmlContent, err := client.FetchHTMLContent(item.URL)
			if err != nil {
				log.Printf("error fetching HTML content for URL %s: %v", item.URL, err)
				continue
			}

			fileURL, err := parser.ExtractFileURLFromHTML(htmlContent)
			if err != nil {
				log.Printf("Error extracting fileUrl for URL %s: %v", item.URL, err)
				continue
			}
			item.FileURL = fileURL
		}

		allResponses = append(allResponses, *response)
		if pageIndex >= response.ResearchItems.TotalPages-1 {
			break
		}
	}

	formattedResponse, err := json.MarshalIndent(allResponses, "", " ")
	if err != nil {
		log.Fatalf("Failed to format response: %v", err)
	}

	if err := os.WriteFile("research.json", formattedResponse, 0644); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Println("Response saved to research.json")
}
