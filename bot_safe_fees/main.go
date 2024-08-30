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

type SafeFeesDetails struct {
	Provider       string            `json:"provider"`
	Fees           FeesSection       `json:"fees"`
	OtherFees      OtherFeesSection  `json:"other_fees"`
	AdditionalInfo AdditionalSection `json:"additional_info"`
}

type FeesSection struct {
	EntranceFee               []FeeItem `json:"entrance_fee"`
	SafeBoxSizeLessThan1000   []FeeItem `json:"safe_box_size_less_than_1000"`
	SafeBoxSize1000To2000     []FeeItem `json:"safe_box_size_1000_to_2000"`
	SafeBoxSize2000To3000     []FeeItem `json:"safe_box_size_2000_to_3000"`
	SafeBoxSizeMoreThan3000   []FeeItem `json:"safe_box_size_more_than_3000"`
	KeyDeposit                FeeItem   `json:"key_deposit"`
	KeyReplacementFee         FeeItem   `json:"key_replacement_fee"`
	SafeDepositBoxDrillingFee FeeItem   `json:"safe_deposit_box_drilling_fee"`
}

type FeeItem struct {
	Text         string    `json:"text"`
	NumericValue *float64  `json:"numeric_value,omitempty"`
	Condition    *string   `json:"condition,omitempty"`
	Amounts      []float64 `json:"amounts"`
}

type OtherFeesSection struct {
	OtherFees *string `json:"other_fees"`
}

type AdditionalSection struct {
	FeeWebsiteLinks string `json:"fee_website_links"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/SafeDepositBoxServiceFee/CompareProductList"
	payloadTemplate := (`{"ProductIdList":"22,50,59,55,44,67,43,21,35,57,28,56,7,61,34","Page":%d,"Limit":3}`)

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95@AdobeOrg=1; AMCV_F915091E62ED182D0A495F95@AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.54; _ctout6672=1; visit_time=13; _ga_NLQFGWVNXN=GS1.1.1725005721.61.1.1725008703.41.0.0`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/OtherFee/SafeDepositBoxServiceFee/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", "wOtMHHCbxVnBmyxoDil8c2CVVE-xPABkb54TQGyS_sTucouzC9oXom7guar3k56sz2tupUOoXZcWMLRLVnBgkTMo92r-UGVvnIRHI4OOjuQ1,d2wxcDXHhYSdYYlu31OZKUYs2IJ4Ig8xWtqVIvH7YsOJK8Bd4hfrA27l98vYRdp4aKPL1AUUxsHmUnz0KZ3zCi2vbuow7BuTpWHRLGye8Sg1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reaading response:", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to parse response body: %v", err)
	}

	totalPages := DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var safeFeesList []SafeFeesDetails
	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}

		req.Header.Set("Accept", "text/plain, */*; q=0.01")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("Cookie", `verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95@AdobeOrg=1; AMCV_F915091E62ED182D0A495F95@AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.54; _ctout6672=1; visit_time=13; _ga_NLQFGWVNXN=GS1.1.1725005721.61.1.1725008703.41.0.0`)
		req.Header.Set("Origin", "https://app.bot.or.th")
		req.Header.Set("Priority", "u=1, i")
		req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/OtherFee/SafeDepositBoxServiceFee/CompareProduct")
		req.Header.Set("Sec-CH-UA", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
		req.Header.Set("Sec-CH-UA-Mobile", "?0")
		req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
		req.Header.Set("VerificationToken", "wOtMHHCbxVnBmyxoDil8c2CVVE-xPABkb54TQGyS_sTucouzC9oXom7guar3k56sz2tupUOoXZcWMLRLVnBgkTMo92r-UGVvnIRHI4OOjuQ1,d2wxcDXHhYSdYYlu31OZKUYs2IJ4Ig8xWtqVIvH7YsOJK8Bd4hfrA27l98vYRdp4aKPL1AUUxsHmUnz0KZ3zCi2vbuow7BuTpWHRLGye8Sg1")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Failed to parse response body: %v", err)
		}

		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			var result SafeFeesDetails

			// Extract the provider name
			result.Provider = cleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())

			// Parse Fees
			result.Fees = parseFees(doc, col)

			// Parse Other Fees
			result.OtherFees = parseOtherFees(doc, col)

			// Parse Additional Info
			result.AdditionalInfo = parseAdditionalInfo(doc, col)

			// Append to the list
			safeFeesList = append(safeFeesList, result)
		}
		time.Sleep(2 * time.Second)
	}

	jsonData, err := json.MarshalIndent(safeFeesList, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile("safe_fees.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Data has been saved to safe_fees.json")
}

func cleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", "")
	text = strings.Join(strings.Fields(text), " ")
	return text
}

func parseFees(doc *goquery.Document, col string) FeesSection {
	return FeesSection{
		EntranceFee:               extractTextArray(doc, fmt.Sprintf("tr.attr-EntranceFee td.%s span", col)),
		SafeBoxSizeLessThan1000:   extractTextArray(doc, fmt.Sprintf("tr.attr-SafeDepositBoxSizeLessThan1000 td.%s span", col)),
		SafeBoxSize1000To2000:     extractTextArray(doc, fmt.Sprintf("tr.attr-SafeDepositBoxSize1000To2000 td.%s span", col)),
		SafeBoxSize2000To3000:     extractTextArray(doc, fmt.Sprintf("tr.attr-SafeDepositBoxSize2000To3000 td.%s span", col)),
		SafeBoxSizeMoreThan3000:   extractTextArray(doc, fmt.Sprintf("tr.attr-SafeDepositBoxSizeMoreThan3000 td.%s span", col)),
		KeyDeposit:                extractFeeItem(doc, fmt.Sprintf("tr.attr-DepositFeeForSafeBoxKey td.%s span", col)),
		KeyReplacementFee:         extractFeeItem(doc, fmt.Sprintf("tr.attr-SafeBoxKeyReplacementFee td.%s span", col)),
		SafeDepositBoxDrillingFee: extractFeeItem(doc, fmt.Sprintf("tr.attr-SafeDepositBoxDrillingFee td.%s span", col)),
	}
}

func extractTextArray(doc *goquery.Document, selector string) []FeeItem {
	var result []FeeItem
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := cleanText(s.Text())
		parts := splitText(text)
		for _, part := range parts {
			entry := FeeItem{
				Text:    part,
				Amounts: extractAmounts(part),
			}
			if num, ok := extractNumericValue(part); ok {
				entry.NumericValue = &num
			}
			if cond := extractCondition(part); cond != "" {
				entry.Condition = &cond
			}
			result = append(result, entry)
		}
	})
	return result
}

func extractFeeItem(doc *goquery.Document, selector string) FeeItem {
	text := cleanText(doc.Find(selector).Text())
	return FeeItem{
		Text:    text,
		Amounts: extractAmounts(text),
	}
}

func extractNumericValue(text string) (float64, bool) {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	match := re.FindString(text)
	if match != "" {
		value, err := strconv.ParseFloat(match, 64)
		if err == nil {
			return value, true
		}
	}
	return 0, false
}

func extractAmounts(text string) []float64 {
	var amounts []float64
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	matches := re.FindAllString(text, -1)
	for _, match := range matches {
		value, err := strconv.ParseFloat(match, 64)
		if err == nil {
			amounts = append(amounts, value)
		}
	}
	return amounts
}

func extractCondition(text string) string {
	if strings.Contains(text, "เงื่อนไข") {
		return text
	}
	return ""
}

func splitText(text string) []string {
	// Split only by "-"
	parts := strings.Split(text, "-")
	var cleanedParts []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			cleanedParts = append(cleanedParts, part)
		}
	}
	return cleanedParts
}

func parseOtherFees(doc *goquery.Document, col string) OtherFeesSection {
	otherFeesText := cleanText(doc.Find(fmt.Sprintf("tr.attr-other td.%s span", col)).Text())
	var otherFees *string
	if otherFeesText == "" {
		otherFees = nil
	} else {
		otherFees = &otherFeesText
	}
	return OtherFeesSection{OtherFees: otherFees}
}

func parseAdditionalInfo(doc *goquery.Document, col string) AdditionalSection {
	link := doc.Find(fmt.Sprintf("tr.attr-Feeurl td.%s a", col)).AttrOr("href", "")
	return AdditionalSection{FeeWebsiteLinks: link}
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
