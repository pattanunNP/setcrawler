package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Define the Go struct to match the JSON structure
type Stock struct {
	Symbol             string  `json:"symbol"`
	LastPrice          float64 `json:"lastPrice"`
	TotalCoverage      int     `json:"totalCoverage"`
	Buy                int     `json:"buy"`
	Hold               int     `json:"hold"`
	Sell               int     `json:"sell"`
	RecommendType      string  `json:"recommendType"`
	MedianTargetPrice  float64 `json:"medianTargetPrice"`
	AverageTargetPrice float64 `json:"averageTargetPrice"`
	Bullish            float64 `json:"bullish"`
	Bearish            float64 `json:"bearish"`
	AnalystHTML        string  `json:"analystHTML"` // New field to store HTML content
}

type Response struct {
	MarketTime string  `json:"marketTime"`
	Overall    []Stock `json:"overall"`
}

func main() {
	// Load the response.json file
	fileName := "response.json"
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Parse the JSON content into the Response struct
	var responseData Response
	err = json.Unmarshal(fileContent, &responseData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Counter for symbols processed
	symbolCount := 0

	// Iterate over each stock symbol and fetch the HTML content
	for i, stock := range responseData.Overall {
		if symbolCount >= 3 {
			break
		}

		// Construct the URL with the stock symbol
		url := fmt.Sprintf("https://www.settrade.com/th/equities/quote/%s/analyst-consensus", stock.Symbol)

		// Create an HTTP GET request
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		// Set necessary headers (adjust as needed)
		req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

		// Send the request
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request for symbol:", stock.Symbol, err)
			continue
		}
		defer res.Body.Close()

		// Read the HTML content
		htmlContent, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body for symbol:", stock.Symbol, err)
			continue
		}

		// Store the HTML content into the struct
		responseData.Overall[i].AnalystHTML = string(htmlContent)

		fmt.Printf("Fetched HTML content for symbol: %s\n", stock.Symbol)

		// Increment the counter
		symbolCount++
	}

	// Save the updated JSON back to a file
	updatedJSON, err := json.MarshalIndent(responseData, "", "  ")
	if err != nil {
		fmt.Println("Error formatting updated JSON:", err)
		return
	}

	err = ioutil.WriteFile("updated_response.json", updatedJSON, 0644)
	if err != nil {
		fmt.Println("Error saving updated JSON to file:", err)
		return
	}

	fmt.Println("Updated JSON saved to updated_response.json")
}
