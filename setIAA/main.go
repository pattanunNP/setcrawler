package main

import (
	"encoding/json"
	"fmt"
	"os"

	"setIAA/pkg/analyst"
	"setIAA/pkg/scraper"
)

func main() {
	// Load the response.json file
	fileName := "response.json"
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	// Parse the JSON content into the Response struct
	var responseData analyst.Response
	err = json.Unmarshal(fileContent, &responseData)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Process each stock symbol and fetch the HTML content
	// symbolCount := 0
	for i := range responseData.Overall {
		// if symbolCount >= 10 {
		// 	break
		// }

		// Construct the URL with the stock symbol
		url := fmt.Sprintf("https://www.settrade.com/th/equities/quote/%s/analyst-consensus", responseData.Overall[i].Symbol)

		// Fetch the HTML content
		htmlContent, err := scraper.FetchHTML(url)
		if err != nil {
			fmt.Printf("Error making request for symbol %s: %v\n", responseData.Overall[i].Symbol, err)
			continue
		}

		// Extract the content inside the specific <script> tag
		scriptContent := scraper.ExtractScriptContent(htmlContent)

		// Extract and assign analyst data
		analystData := analyst.ExtractAnalystData(scriptContent)
		responseData.Overall[i].AnalystData = analystData

		fmt.Printf("Extracted script content for symbol %s\n", responseData.Overall[i].Symbol)

		// symbolCount++
	}

	// Save the updated JSON back to a file
	updatedJSON, err := json.MarshalIndent(responseData, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting updated JSON: %v\n", err)
		return
	}

	err = os.WriteFile("updated_response.json", updatedJSON, 0644)
	if err != nil {
		fmt.Printf("Error saving updated JSON to file: %v\n", err)
		return
	}

	fmt.Println("Updated JSON saved to updated_response.json")
}
