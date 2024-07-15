package MajorShareHolder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchDetailedShareholderData(postURL, cookieStr, formData string) (string, error) {
	req, err := http.NewRequest("POST", postURL, bytes.NewBufferString(formData))
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %w", err)
	}

	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("origin", "https://www.setsmart.com")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("referer", "https://www.setsmart.com/ism/majorshareholder.html?locale=en_US")
	req.Header.Add("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "macOS")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error maing POST requesy: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

// ExtractOptionValues extracts the values of all <option> tags within the specified select box
func ExtractOptionValues(htmlStr, selectName string) ([]string, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, fmt.Errorf("error loading HTML document: %w", err)
	}

	var values []string

	// Find the select box with the specified name
	doc.Find(fmt.Sprintf("select[name=%s] option", selectName)).Each(func(index int, item *goquery.Selection) {
		value, exists := item.Attr("value")
		if exists {
			values = append(values, value)
		}
	})

	return values, nil
}

// FetchDataForOptionValues sends a POST request for each option value and returns the results
func FetchDataForOptionValues(postURL, cookieStr string, optionValues []string) ([]string, error) {
	var results []string

	for _, value := range optionValues {
		formData := fmt.Sprintf("radChoice=1&txtSymbol=SCB&radShow=2&hidAction=&hidLastContentType=&lstFreeFloatDate=%s", value)
		response, err := FetchDetailedShareholderData(postURL, cookieStr, formData)
		if err != nil {
			return nil, fmt.Errorf("error fetching data for option value %s: %w", value, err)
		}
		results = append(results, response)
	}

	return results, nil
}

// ExtractDataFromTable extracts data from <td> elements of a specific table
func ExtractDataFromTable(doc *goquery.Document, tableIndex int) (map[string]string, error) {
	data := make(map[string]string)

	// Find the specific table with the class "tfont" by index
	table := doc.Find("table.tfont").Eq(tableIndex)
	if table.Length() == 0 {
		return nil, fmt.Errorf("table with class 'tfont' not found at index %d", tableIndex)
	}

	// Iterate over the <td> elements
	table.Find("tr").Each(func(index int, row *goquery.Selection) {
		name := row.Find("td.table-bold").Text()
		value := row.Find("td.table").Text()
		if name != "" && value != "" {
			// Remove whitespace and newlines
			name = strings.TrimSpace(name)
			value = strings.TrimSpace(value)
			data[name] = value
		}
	})

	return data, nil
}

// ExtractTableRows extracts text content from all rows in a table
func ExtractTableRows(doc *goquery.Document, tableIndex int) ([]string, error) {
	var rows []string

	// Find the specific table with the class "tfont" by index
	table := doc.Find("table.tfont").Eq(tableIndex)
	if table.Length() == 0 {
		return nil, fmt.Errorf("table with class 'tfont' not found at index %d", tableIndex)
	}

	// Iterate over the rows
	table.Find("tr").Each(func(index int, row *goquery.Selection) {
		rows = append(rows, row.Text())
	})

	return rows, nil
}

// SaveDataAsJSON saves the provided data map as a JSON file
func SaveDataAsJSON(data map[string]map[string]string, jsonFilename string) error {
	// Convert the data to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Save the JSON data to a file
	err = os.WriteFile(jsonFilename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	return nil
}
