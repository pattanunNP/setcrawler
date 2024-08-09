package main

import (
	"fmt"
	"log"
	"setResearch/httpclient"
	"setResearch/research"
	"setResearch/utils"
)

func main() {
	client := httpclient.NewClient()

	var allItems []research.ResearchItem
	pageIndex := 0

	for {
		responseData, err := research.FetchResearchItems(client, pageIndex)
		if err != nil {
			log.Fatalf("Error fetching research items: %v", err)
		}

		for i := range responseData.ResearchItems.Items {
			item := &responseData.ResearchItems.Items[i]
			if item.URL != "" {
				fmt.Printf("Processing item: %s\n", item.Title)
				err := research.ProcessResearchItem(client, item)
				if err != nil {
					log.Printf("Error processing item %s: %v\n", item.Title, err)
					continue
				}
				allItems = append(allItems, *item)
			}
		}

		if responseData.ResearchItems.PageIndex >= responseData.ResearchItems.TotalPages-1 {
			break
		}

		pageIndex++

	}

	fileName := "response_with_pdf_content.json"
	err := utils.SaveToFile(fileName, allItems)
	if err != nil {
		log.Fatalf("Error saving response to file: %v", err)
	}
	fmt.Printf("response with PDF content saved to %s\n", fileName)
}
