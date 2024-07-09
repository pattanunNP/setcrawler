package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"login_token/utils"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// FetchJSON makes a GET request to the provided URL and returns the response body.
func FetchJSON(apiURL string) ([]byte, error) {
	client, cookies, err := utils.SetupCookieAndToken()
	if err != nil {
		return nil, fmt.Errorf("error setting up cookies and token: %w", err)
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add cookies to the request
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// ExtractSymbols extracts the "symbol" fields from the provided JSON content.
func ExtractSymbols(jsonContent []byte) ([]string, error) {
	// Parse the JSON content assuming the structure is an object with a "stocks" key containing the array
	var jsonResponse map[string]interface{}
	err := json.Unmarshal(jsonContent, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Extract the "symbol" fields
	var symbols []string
	if records, ok := jsonResponse["stocks"].([]interface{}); ok {
		for _, record := range records {
			if recordMap, ok := record.(map[string]interface{}); ok {
				if symbol, ok := recordMap["symbol"].(string); ok {
					symbols = append(symbols, symbol)
				}
			}
		}
	}

	return symbols, nil
}

func ReadSymbolsFromFile(fileName string) ([]string, error) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	symbols, err := ExtractSymbols(fileContent)
	if err != nil {
		return nil, fmt.Errorf("error extracting symbols: %w", err)
	}

	return symbols, nil
}

// MakePostRequests makes POST requests for each symbol and each locale.
func MakePostRequests(symbols []string, locales []string) error {
	client, cookies, err := utils.SetupCookieAndToken()
	if err != nil {
		return fmt.Errorf("error setting up cookies and token: %w", err)
	}

	apiURL := "https://www.setsmart.com/ism/companyprofile.html"

	for _, symbol := range symbols {
		for _, locale := range locales {
			// Form data
			form := url.Values{}
			form.Add("symbol", symbol)
			form.Add("locale", locale)

			req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(form.Encode()))
			if err != nil {
				return fmt.Errorf("error creating request: %w", err)
			}

			// Add cookies to the request
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}

			// Add headers
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
			req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
			req.Header.Set("Cache-Control", "max-age=0")
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Origin", "https://www.setsmart.com")
			req.Header.Set("Referer", "https://www.setsmart.com/ism/companyprofile.html?locale="+locale)
			req.Header.Set("Sec-CH-UA", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
			req.Header.Set("Sec-CH-UA-Mobile", "?0")
			req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
			req.Header.Set("Sec-Fetch-Dest", "document")
			req.Header.Set("Sec-Fetch-Mode", "navigate")
			req.Header.Set("Sec-Fetch-Site", "same-origin")
			req.Header.Set("Sec-Fetch-User", "?1")
			req.Header.Set("Upgrade-Insecure-Requests", "1")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

			// Debugging: Print request details
			fmt.Printf("Making request for symbol: %s, locale: %s\n", symbol, locale)
			fmt.Printf("Request URL: %s\n", req.URL.String())
			fmt.Printf("Request Body: %s\n", form.Encode())
			fmt.Printf("Request Headers: %v\n", req.Header)
			fmt.Printf("Request Cookies: %v\n", req.Cookies())

			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error sending request: %w", err)
			}
			defer resp.Body.Close()

			// Debugging: Print response status
			fmt.Printf("Response Status: %s\n", resp.Status)

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("error reading response body: %w", err)
			}

			// Debugging: Print response details
			fmt.Printf("Response for symbol %s with locale %s:\n%s\n", symbol, locale, string(body))

			// Check if session timed out or permission denied
			if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
				fmt.Printf("Permission denied or session timed out for symbol %s with locale %s\n", symbol, locale)
			}
		}
	}

	return nil
}

// UpdateAndSaveSymbolsStatus fetches data from both URLs, compares the symbols, updates the status, and saves the result to a file.
func UpdateAndSaveSymbolsStatus(allAPIURL string, activeAPIURL string, outputFileName string) error {
	allJSON, err := FetchJSON(allAPIURL)
	if err != nil {
		return fmt.Errorf("error fetching all stocks: %w", err)
	}

	activeJSON, err := FetchJSON(activeAPIURL)
	if err != nil {
		return fmt.Errorf("error fetching active stocks: %w", err)
	}

	activeSymbols, err := ExtractSymbols(activeJSON)
	if err != nil {
		return fmt.Errorf("error extracting active symbols: %w", err)
	}

	activeSymbolsMap := make(map[string]bool)
	for _, symbol := range activeSymbols {
		activeSymbolsMap[symbol] = true
	}

	var allStocksResponse map[string]interface{}
	err = json.Unmarshal(allJSON, &allStocksResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling all stocks JSON: %w", err)
	}

	if records, ok := allStocksResponse["stocks"].([]interface{}); ok {
		for _, record := range records {
			if recordMap, ok := record.(map[string]interface{}); ok {
				if symbol, ok := recordMap["symbol"].(string); ok {
					if activeSymbolsMap[symbol] {
						recordMap["status"] = "active"
					} else {
						recordMap["status"] = "inactive"
					}
				}
			}
		}
	}

	updatedJSON, err := json.MarshalIndent(allStocksResponse, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling updated JSON: %w", err)
	}
	fmt.Println(string(updatedJSON))

	outputDir := filepath.Dir(outputFileName)
	if outputDir != "." {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
	}

	err = os.WriteFile(outputFileName, updatedJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing updated JSON to file: %w", err)
	}

	fmt.Printf("Updated response saved to %s\n", outputFileName)
	return nil
}
