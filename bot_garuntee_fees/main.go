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
	"time"

	"github.com/PuerkitoBio/goquery"
)

type GarunteeDetails struct {
	Provider       string            `json:"provider"`
	Fees           FeeDetails        `json:"fees"`
	OtherFees      OtherFeeDetails   `json:"other_fees"`
	AdditionalInfo AdditionalDetails `json:"additional_info"`
}

type FeeDetails struct {
	Merchandise      FeeItem `json:"merchandise"`
	Advance          FeeItem `json:"advance"`
	Borrowing        FeeItem `json:"borrowing"`
	Bid              FeeItem `json:"bid"`
	Performance      FeeItem `json:"performance"`
	Retention        FeeItem `json:"retention"`
	ElectricityWater FeeItem `json:"electricity_water"`
	Tax              FeeItem `json:"tax"`
	Others           FeeItem `json:"others"`
}

type FeeItem struct {
	OriginalText  string  `json:"original_text"`
	MaxPercentage float64 `json:"max_percentage,omitempty"`
	MinFee        float64 `json:"min_fee,omitempty"`
	Conditions    string  `json:"conditions,omitempty"`
}

type OtherFeeDetails struct {
	Other *string `json:"other,omitempty"`
}

type AdditionalDetails struct {
	FeeWebsiteLink *string `json:"fee_website_link,omitempty"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/GuaranteeIssuingServiceFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"27,60,45,57,58,9,38,7,47,53,54,40,39,56,30,5,4,61,16,11,55,26,43,41,21,34,52,20,8,33","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Set the headers
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/OtherFee/GuaranteeIssuingServiceFee/CompareProduct")
	req.Header.Set("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", "KjpDbi0v7Ntdzk9CczUA6AkY5jCBbdSGms1jeO7dv51-PW3TVMuU_eJk4u5TPXzmKqwDKsKiL_TNk8tMM1NzBtn3uv8eHwmwDuRfQvz-0rs1,DcNTb8wy_TOW9C0rKE_JgvartfWEgxYhWCOVkoT9Lg4fA19fcoXfRVCSSUaFF2t3uHDgS_D6pJrooT2nPIGQ2lK9OdQR_tuO54TXJftG58Y1")
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to parse response body: %v", err)
	}

	totalPages := DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var garunteeDetailsList []GarunteeDetails

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}

		// Set the headers
		req.Header.Set("accept", "text/plain, */*; q=0.01")
		req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
		req.Header.Set("content-type", "application/json; charset=UTF-8")
		req.Header.Set("origin", "https://app.bot.or.th")
		req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/OtherFee/GuaranteeIssuingServiceFee/CompareProduct")
		req.Header.Set("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"macOS"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
		req.Header.Set("verificationtoken", "KjpDbi0v7Ntdzk9CczUA6AkY5jCBbdSGms1jeO7dv51-PW3TVMuU_eJk4u5TPXzmKqwDKsKiL_TNk8tMM1NzBtn3uv8eHwmwDuRfQvz-0rs1,DcNTb8wy_TOW9C0rKE_JgvartfWEgxYhWCOVkoT9Lg4fA19fcoXfRVCSSUaFF2t3uHDgS_D6pJrooT2nPIGQ2lK9OdQR_tuO54TXJftG58Y1")
		req.Header.Set("x-requested-with", "XMLHttpRequest")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Failed to parse response body: %v", err)
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			var result GarunteeDetails

			// Extract the provider name
			result.Provider = cleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())

			// List of fee categories to iterate through
			feeCategories := []struct {
				attr  string
				field *FeeItem
			}{
				{"attr-Merchandise", &result.Fees.Merchandise},
				{"attr-Advance", &result.Fees.Advance},
				{"attr-Borrowing", &result.Fees.Borrowing},
				{"attr-Bid", &result.Fees.Bid},
				{"attr-Performance", &result.Fees.Performance},
				{"attr-Retention", &result.Fees.Retention},
				{"attr-ElectricityWater", &result.Fees.ElectricityWater},
				{"attr-Tax", &result.Fees.Tax},
				{"attr-Others", &result.Fees.Others},
			}

			// Extract fee details for each category
			for _, category := range feeCategories {
				originalText := cleanText(doc.Find(fmt.Sprintf("tr.%s td.%s span", category.attr, col)).Text())
				maxPercent, minFee, conditions := parseFeeText(originalText)
				*category.field = FeeItem{
					OriginalText:  originalText,
					MaxPercentage: maxPercent,
					MinFee:        minFee,
					Conditions:    conditions,
				}
			}

			// Extract other fees for this provider
			otherFeeText := cleanText(doc.Find(fmt.Sprintf("tr.attr-header.attr-other ~ tr td.%s span", col)).Text())
			if otherFeeText != "" {
				result.OtherFees.Other = &otherFeeText
			} else {
				result.OtherFees.Other = nil
			}

			// Extract additional information for this provider
			additionalLink := doc.Find(fmt.Sprintf("tr.attr-header.attr-additional ~ tr td.%s a.prod-url", col)).AttrOr("href", "")
			if additionalLink != "" {
				result.AdditionalInfo.FeeWebsiteLink = &additionalLink
			} else {
				result.AdditionalInfo.FeeWebsiteLink = nil
			}

			// Append this result to the garuntee details list
			garunteeDetailsList = append(garunteeDetailsList, result)
		}
		time.Sleep(2 * time.Second)
	}

	jsonData, err := json.MarshalIndent(garunteeDetailsList, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile("garuntee_fees.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Data has been saved to garuntee_fees.json")
}

func parseFeeText(text string) (float64, float64, string) {
	var maxPercentage, minFee float64
	var conditions string

	// Regex to find percentage and minimum fee
	rePercent := regexp.MustCompile(`(\d+(\.\d+)?)%`)
	reMinFee := regexp.MustCompile(`ขั้นต่ำ (\d+)( บาท| USD)`)

	// Extract percentage
	percentMatch := rePercent.FindStringSubmatch(text)
	if len(percentMatch) > 1 {
		maxPercentage, _ = strconv.ParseFloat(percentMatch[1], 64)
	}

	// Extract minimum fee
	minFeeMatch := reMinFee.FindStringSubmatch(text)
	if len(minFeeMatch) > 1 {
		minFee, _ = strconv.ParseFloat(minFeeMatch[1], 64)
	}

	// Extract conditions (anything after 'เงื่อนไข:')
	conditionIndex := strings.Index(text, "เงื่อนไข:")
	if conditionIndex != -1 {
		conditions = strings.TrimSpace(text[conditionIndex+len("เงื่อนไข:"):])
	}

	return maxPercentage, minFee, conditions
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
