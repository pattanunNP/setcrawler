package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Fee struct {
	Description *string  `json:"description"`
	Conditions  []string `json:"conditions,omitempty"`
}

type Fees struct {
	TransferNextDay           Fee `json:"transfer_next_day"`
	TransferSameDay100K       Fee `json:"transfer_same_day_100k"`
	TransferSameDay500K       Fee `json:"transfer_same_day_500k"`
	TransferSameDayTwoMillion Fee `json:"transfer_same_day_two_million"`
	PromptPay100K             Fee `json:"promptpay_100k"`
	PromptPayTwoMillion       Fee `json:"promptpay_two_million"`
	DirectCreditInBranch      Fee `json:"direct_credit_in_branch"`
	DirectCreditAcrossBranch  Fee `json:"direct_credit_across_branch"`
	DirectDebitInBranch       Fee `json:"direct_debit_in_branch"`
	DirectDebitAcrossBranch   Fee `json:"direct_debit_across_branch"`
	Other                     Fee `json:"other_fees"`
}

type AdditionalInfo struct {
	FeeWebsite *string `json:"fee_website"`
}

type BulkPayment struct {
	Provider       string         `json:"provider"`
	Fees           Fees           `json:"Fees"`
	OtherFees      Fee            `json:"other_fees"`
	AdditionalInfo AdditionalInfo `json:"additional_info"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/BulkPaymentFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"12,34,15,10,1,44,22,14,33,39,9,2,8,31,40,32,38,13,6,36,11,43,16,42,24,3,28","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	addHeaders(req)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Error loading HTTP response body: %v", err)
	}

	totalPages := determineTotalPage(doc)
	if totalPages == 0 {
		log.Fatalf("Could not determine the total number of pages")
	}

	var bulkPayments []BulkPayment

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
			provider := strings.TrimSpace(doc.Find("th.attr-header.attr-prod.font-black.text-center.cmpr-col." + col + " span").Last().Text())

			// Create a new BulkPayment for each provider
			bulkPayment := BulkPayment{
				Provider: provider,
				Fees: Fees{
					TransferNextDay:           extractFee(doc, "tr.attr-NextDay", col),
					TransferSameDay100K:       extractFee(doc, "tr.attr-SameDayOneHundredK", col),
					TransferSameDay500K:       extractFee(doc, "tr.attr-SameDayFiveHundredK", col),
					TransferSameDayTwoMillion: extractFee(doc, "tr.attr-SameDayTwoMillion", col),
					PromptPay100K:             extractFee(doc, "tr.attr-PromptPayOneHundredK", col),
					PromptPayTwoMillion:       extractFee(doc, "tr.attr-PromptPayTwoMillion", col),
					DirectCreditInBranch:      extractFee(doc, "tr.attr-DirectCreditInbranch", col),
					DirectCreditAcrossBranch:  extractFee(doc, "tr.attr-DirectCreditAccross", col),
					DirectDebitInBranch:       extractFee(doc, "tr.attr-DirectDebitInbranch", col),
					DirectDebitAcrossBranch:   extractFee(doc, "tr.attr-DirectDebitAccross", col),
				},
				OtherFees: extractFee(doc, "tr.attr-other", col),
			}

			// Extract additional info (website links)
			website := doc.Find("tr.attr-Feeurl td.cmpr-col."+col+" a.prod-url").AttrOr("href", "")
			if website != "" {
				bulkPayment.AdditionalInfo.FeeWebsite = &website
			}

			bulkPayments = append(bulkPayments, bulkPayment)
		}
	}

	jsonData, err := json.MarshalIndent(bulkPayments, "", " ")
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	err = os.WriteFile("bulk_payment.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Printf("Data Successfully Saving to bulk_payment.json")
}

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

func addHeaders(req *http.Request) {
	// Set headers
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", "verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95@AdobeOrg=1; AMCV_F915091E62ED182D0A495F95@AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT=\"z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0\"; _uid6672=16B5DEBD.36; _ctout6672=1; visit_time=9; _ga_NLQFGWVNXN=GS1.1.1724760674.44.1.1724760693.41.0.0")
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/BulkPaymentFee/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", "bzpgEYIrRyCs_becKiSUxmAAPN2EeBZ0oJIQGTydWdaE8wKPw_qUPtxZxrbsp_7l8NQ2d7XrVsLVI_aEaEYSTR9v4k9geAh02j8qJFq_JOI1,Jr8T4mleZa0qNVizfWUEHSvmjZzMPhPSuETOd-MKvrVme3fxe6ePouczaiCanh5W9UV6b8pxPb-Iy2VWcnsbOWe8YuxLjHZCQsECjlydAEo1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

// extractFee extracts fee information based on the column and row selectors
func extractFee(doc *goquery.Document, rowSelector, col string) Fee {
	var fee Fee
	doc.Find(rowSelector + " td.cmpr-col." + col).Each(func(i int, s *goquery.Selection) {
		description := strings.TrimSpace(s.Find("span").First().Text())
		if description != "" {
			fee.Description = &description
		}

		conditionText := s.Find("div.collapse span.text-primary").Text()
		if conditionText != "" {
			// Remove newline characters and clean up spaces
			conditionText = strings.ReplaceAll(conditionText, "\n", " ")
			conditionText = strings.TrimSpace(conditionText)

			// Use strings.Fields to split by any whitespace and then join with a single space
			cleanedText := strings.Join(strings.Fields(conditionText), " ")

			conditions := strings.Split(cleanedText, "-")
			for _, cond := range conditions {
				cleanedCond := strings.TrimSpace(cond)
				if cleanedCond != "" {
					fee.Conditions = append(fee.Conditions, cleanedCond)
				}
			}
		}
	})
	return fee
}
