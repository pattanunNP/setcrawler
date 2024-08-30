package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Define the structure of JSON output using nested structs
type CustodianDetails struct {
	Provider      string          `json:"provider"`
	Fees          []Fees          `json:"fees"`
	OtherFees     OtherFeeDetails `json:"other_fees"`
	AdditionalURL AdditionalInfo  `json:"additional_info"`
}

type Fees struct {
	Min          float64 `json:"min,omitempty"`
	Max          float64 `json:"max,omitempty"`
	OriginalText string  `json:"original_text"`
	Description  string  `json:"description,omitempty"`
	Condition    string  `json:"condition,omitempty"`
}

type OtherFeeDetails struct {
	Description *string `json:"description"`
}

type AdditionalInfo struct {
	Website *string `json:"website"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/CustodianServiceFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"4,41,30,59,20,27,21,48,57,53,39,54","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/OtherFee/CustodianServiceFee/CompareProduct")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("VerificationToken", "e700mfmLRr_25GqxWWHnm1vAFGdmxHVXEYG-4iGf5LqLdD6kc8xgn6lK4WfCQ1RQXVXZ8qln5-VdgJ8KRzdUOcQSVpaV8upjA2_kPY57JG81,eE4EtMdzu-er2AIvdoCxC36djIIwneO-Eh6MEFh71GbHdpTiklEC8-ti7ZzOVclr6Geat1wTNzuJV1btFFIaY3NgrsQXppp8fu_zkwfGLho1")
	req.Header.Set("Cookie", "verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error loading HTML document:", err)
		return
	}

	totalPages := DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var custodianDetailsList []CustodianDetails
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}

		// Set headers
		req.Header.Set("Accept", "text/plain, */*; q=0.01")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("Origin", "https://app.bot.or.th")
		req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/OtherFee/CustodianServiceFee/CompareProduct")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("VerificationToken", "e700mfmLRr_25GqxWWHnm1vAFGdmxHVXEYG-4iGf5LqLdD6kc8xgn6lK4WfCQ1RQXVXZ8qln5-VdgJ8KRzdUOcQSVpaV8upjA2_kPY57JG81,eE4EtMdzu-er2AIvdoCxC36djIIwneO-Eh6MEFh71GbHdpTiklEC8-ti7ZzOVclr6Geat1wTNzuJV1btFFIaY3NgrsQXppp8fu_zkwfGLho1")
		req.Header.Set("Cookie", "verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=")

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			fmt.Println("Error loading HTML document:", err)
			return
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			var result CustodianDetails

			// Extract the provider name
			result.Provider = cleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())

			// Extract custodian fee details for this provider
			var feesList []Fees
			doc.Find(fmt.Sprintf("tr.attr-header.attr-fee ~ tr td.%s span", col)).Each(func(i int, s *goquery.Selection) {
				text := cleanText(s.Text())
				if text != "" {
					fee := Fees{
						OriginalText: text,
					}

					// Extract min and max values from the fee description
					min, max := extractMinMax(text)
					if min != 0 || max != 0 {
						fee.Min = min
						fee.Max = max
					} else {
						// If not a range, treat as condition
						fee.Condition = text
					}

					feesList = append(feesList, fee)
				}
			})
			result.Fees = feesList

			// Extract other fees for this provider
			var otherFeeDesc *string
			otherFeeText := cleanText(doc.Find(fmt.Sprintf("tr.attr-header.attr-other ~ tr td.%s span", col)).Text())
			if otherFeeText != "" {
				otherFeeDesc = &otherFeeText
			}
			result.OtherFees.Description = otherFeeDesc

			// Extract additional information for this provider
			var additionalURL *string
			additionalLink := doc.Find(fmt.Sprintf("tr.attr-header.attr-additional ~ tr td.%s a.prod-url", col)).AttrOr("href", "")
			if additionalLink != "" {
				additionalURL = &additionalLink
			}
			result.AdditionalURL.Website = additionalURL

			// Append this result to the custodian details list
			custodianDetailsList = append(custodianDetailsList, result)
		}
	}

	jsonData, err := json.MarshalIndent(custodianDetailsList, "", " ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	err = os.WriteFile("custodian.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("Data successfully written to custodian.json")
}

func DetermineTotalPage(doc *goquery.Document) int {
	totalPages := 1

	doc.Find("ul.pagination li a").Each(func(i int, s *goquery.Selection) {
		pagenum, exists := s.Attr("data-page")
		if exists {
			page, err := strconv.Atoi(pagenum)
			if err == nil && page > totalPages {
				totalPages = page
			}
		}
	})
	return totalPages
}

func cleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", "")
	text = strings.Join(strings.Fields(text), " ")
	return text
}

// Extract min and max values from the given text
func extractMinMax(text string) (float64, float64) {
	re := regexp.MustCompile(`(\d+(\.\d+)?)`)
	matches := re.FindAllString(text, -1)

	if len(matches) >= 2 {
		min, _ := strconv.ParseFloat(matches[0], 64)
		max, _ := strconv.ParseFloat(matches[1], 64)
		return min, max
	} else if len(matches) == 1 {
		value, _ := strconv.ParseFloat(matches[0], 64)
		return value, value
	}
	return 0, 0
}
