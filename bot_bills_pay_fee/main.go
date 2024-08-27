package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// BillsPaymentFees represents the structure for the JSON output
type BillsPaymentFees struct {
	Provider       string         `json:"provider"`
	Fees           FeeCategories  `json:"fees"`
	AdditionalInfo AdditionalInfo `json:"additional_info,omitempty"`
}

// FeeCategories holds various fee categories
type FeeCategories struct {
	ElectricityBill            []FeeDetail `json:"electricity_bill,omitempty"`
	WaterBill                  []FeeDetail `json:"water_bill,omitempty"`
	PhoneOrInternetBill        []FeeDetail `json:"phone_or_internet_bill,omitempty"`
	InsuranceBill              []FeeDetail `json:"insurance_bill,omitempty"`
	VehicleRegistrationBill    []FeeDetail `json:"vehicle_registration_bill,omitempty"`
	MotorcycleRegistrationBill []FeeDetail `json:"motorcycle_registration_bill,omitempty"`
	TaxBill                    []FeeDetail `json:"tax_bill,omitempty"`
	OtherUtilitiesBill         []FeeDetail `json:"other_utilities_bill,omitempty"`
	ProductOrServiceBill       []FeeDetail `json:"product_or_service_bill,omitempty"`
}

// FeeDetail represents a single fee detail
type FeeDetail struct {
	Description string   `json:"description"`
	Details     []string `json:"details"`
}

// AdditionalInfo represents additional information such as URLs
type AdditionalInfo struct {
	FeeURL string `json:"fee_url"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/BillPaymentFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"2-0785800001,2-0785800005,5-0785800005,17-0785800004,17-0785800005,157479-0785800005,26-0785800004,26-0785800005,26-0785800008,150920-0785800001,194031-0785800001,9-0785800003,9-0785800004,9-0785800005,9-0785800007,162152-0785800001,162152-0785800002,5-0785800002,5-0785800003,5-0785800004,241-0785800001,2-0785800004,15-0785800001,13518965-0785800003,163579-0785800003,9-0785800002,237-0785800005,15-0785800005,163579-0785800007,13518965-0785800007,157479-0785800004,471989-0785800005,162568-0785800001,2-0785800002,10651-0785800002,10651-0785800005,471989-0785800002,163579-0785800002,13518965-0785800002,6-0785800002,6-0785800003,6-0785800004,6-0785800005,13519357-0785800002,13519357-0785800003,237-0785800002,237-0785800004,241-0785800002,162568-0785800002,4-0785800004,4-0785800005,15-0785800002,157479-0785800002,13519357-0785800005,4-0785800002,17-0785800002,17-0785800003,17-0785800001,13518965-0785800005,163579-0785800005,5-0785800008,237-0785800001,2-0785800003,5-0785800006,5-0785800007,4-0785800003,1036847-0785800005,422974-0785800007,15-0785800004,15-0785800007,870537-0785800005,163579-0785800001,13518965-0785800001,150920-0785800002,150920-0785800003,150920-0785800005,194031-0785800002,517619-0785800003,162152-0785800004,162152-0785800005,37-0785800007,892849-0785800005,16-0785800001,16-0785800002,16-0785800003,16-0785800004,16-0785800005,162568-0785800004,162568-0785800005,237-0785800003,163222-0785800004,163222-0785800005,390797-0785800002,390797-0785800004,33-0785800007","Page":%d,"Limit":3}`

	// Initial request to determine the number of pages
	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add headers to the request
	addHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	totalPages := determineTotalPage(doc)
	if totalPages == 0 {
		fmt.Println("Could not determine the total number of pages")
		return
	}

	var billsPaymentFees []BillsPaymentFees

	// Loop through each page to gather all data
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			fmt.Println("Error creating request for page", page, ":", err)
			continue
		}

		// Set headers
		addHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request for page", page, ":", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request for page %d failed with status: %d\n", page, resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body for page", page, ":", err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			fmt.Println("Error parsing HTML for page", page, ":", err)
			continue
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			provider := doc.Find("th.attr-header.attr-prod.font-black.text-center.cmpr-col." + col + " span").Last().Text()

			var feeURL string
			doc.Find("tr.attr-header.attr-Feeurl").Each(func(index int, row *goquery.Selection) {
				row.Find("td.cmpr-col." + col).Each(func(index int, cell *goquery.Selection) {
					feeURL = cell.Find("a.prod-url").AttrOr("href", "")
				})
			})

			billsPaymentFees = append(billsPaymentFees, BillsPaymentFees{
				Provider: provider,
				Fees:     extractFeeCategories(doc, col),
				AdditionalInfo: AdditionalInfo{
					FeeURL: feeURL,
				},
			})
		}
	}

	jsonData, err := json.MarshalIndent(billsPaymentFees, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	err = os.WriteFile("bills_payment_fee.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to write JSON to file:", err)
		return
	}

	fmt.Println("Data successfully saved to bills_payment_fee.json")
}

// determineTotalPage determines the total number of pages from the pagination element
func determineTotalPage(doc *goquery.Document) int {
	totalPages := 1

	doc.Find("ul.pagination li a").Each(func(i int, s *goquery.Selection) {
		pageNum, exists := s.Attr("data-page")
		if exists {
			page, err := strconv.Atoi(pageNum)
			if err == nil && page > totalPages {
				totalPages = page
			}
		}
	})
	return totalPages
}

// addHeaders adds necessary headers to the HTTP request
func addHeaders(req *http.Request) {
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.35; _ctout6672=1; _ga_NLQFGWVNXN=GS1.1.1724755244.43.1.1724755258.46.0.0`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/BillPaymentFee/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", "_PIlaRXOFAdPaxxKZ3M3jis3-452v3ifiSSjLwk_qX0fw6YBOcyRwRRQ5-iLDsJupX5VbnkiKnZs8Z-G3VIKT5tD4OtRpKSNPRuLwYQwsnU1,J2csEOu_X6ojPPFpD8SCeCvoviJ333MfxjFG0SHcfuOUkvNGRI8pNAmjwSIkle9LpbQMDHyzAqEGmgRilrtFOUS5vBEbxPpRB4dTHmNOTQM1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

// extractFeeCategories extracts fee details for each provider based on the column class
func extractFeeCategories(doc *goquery.Document, col string) FeeCategories {
	var fees FeeCategories

	fees.ElectricityBill = extractFeeDetails(doc, col, "ElectricityBill")
	fees.WaterBill = extractFeeDetails(doc, col, "WaterBill")
	fees.PhoneOrInternetBill = extractFeeDetails(doc, col, "PhoneOrInternetBill")
	fees.InsuranceBill = extractFeeDetails(doc, col, "InsuranceBill")
	fees.VehicleRegistrationBill = extractFeeDetails(doc, col, "VehicleRegistrationBill")
	fees.MotorcycleRegistrationBill = extractFeeDetails(doc, col, "MotorcycleRegistrationBill")
	fees.TaxBill = extractFeeDetails(doc, col, "TaxBill")
	fees.OtherUtilitiesBill = extractFeeDetails(doc, col, "OtherUtilitiesBill")
	fees.ProductOrServiceBill = extractFeeDetails(doc, col, "ProductOrServiceBill")

	return fees
}

// extractFeeDetails extracts the fee details for a specific fee type and column
func extractFeeDetails(doc *goquery.Document, col, feeType string) []FeeDetail {
	var details []FeeDetail
	doc.Find("tr.attr-header.attr-" + feeType + " td.cmpr-col." + col).Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		lines := strings.Split(text, "\n")
		cleanedLines := cleanTextLines(lines)

		if len(cleanedLines) > 0 {
			detail := FeeDetail{
				Description: feeType,
				Details:     cleanedLines,
			}
			details = append(details, detail)
		}
	})
	return details
}

// cleanTextLines trims whitespace and removes empty lines from the slice of strings
func cleanTextLines(lines []string) []string {
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}
