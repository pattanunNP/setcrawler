package main

import (
	"encoding/json"
	"fmt"
	"log"
	"nano_finance/pkg/httpclient"
	model "nano_finance/pkg/models"
	"nano_finance/pkg/parser"
	"os"
)

func main() {
	// Fetch the first page to determine the total number of pages
	initialPayload := `{"ProductIdList":"192,50,52,236,133,180,139,137,141,59,256,182,257,49,181,32,258,135,242,243,248,132,145,33,218,82,85,240,233,239,237,92,51,194,191,41,249,206,96,245,56,229,10,231,11,196","Page":1,"Limit":3}`
	firstPageBody, err := httpclient.FetchData(initialPayload)
	if err != nil {
		log.Fatalf("Error fetching first page: %v", err)
	}

	// Determine the total number of pages
	totalPages := parser.DetermineTotalPages(firstPageBody)
	if totalPages == 0 {
		log.Fatalf("Could not determine the total number of pages")
	}

	var allProducts []model.ProductInfo

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(`{"ProductIdList":"192,50,52,236,133,180,139,137,141,59,256,182,257,49,181,32,258,135,242,243,248,132,145,33,218,82,85,240,233,239,237,92,51,194,191,41,249,206,96,245,56,229,10,231,11,196","Page":%d,"Limit":3}`, page)
		pageBody, err := httpclient.FetchData(payload)
		if err != nil {
			log.Printf("Error fetching data for page %d: %v", page, err)
			continue
		}

		products, err := parser.ParseProductData(pageBody)
		if err != nil {
			log.Printf("Error parsing product data for page %d: %v", page, err)
			continue
		}
		allProducts = append(allProducts, products...)
	}

	resultJSON, err := json.MarshalIndent(allProducts, "", " ")
	if err != nil {
		log.Fatalf("Error marshling to JSON: %v", err)
	}

	err = os.WriteFile("nanofinace.json", resultJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
	fmt.Println("Extracted data saved to nanofinance.json")

}
