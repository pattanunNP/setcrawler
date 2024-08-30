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

type FinancialStatus struct {
	FinancialCertificationForVisa       Certification `json:"financial_certification_for_visa"`
	FinancialCertificationForGovernment Certification `json:"financial_certification_for_government"`
	FinancialCertificationForAudit      Certification `json:"financial_certification_for_audit"`
}

type Certification struct {
	OriginalText string `json:"original_text"`
	Value        int    `json:"value"`
}

type CreditCertification struct {
	ConditionalCreditCertification   []Certification `json:"conditional_credit_certification"`
	UnconditionalCreditCertification []Certification `json:"unconditional_credit_certification"`
}

type OtherFees struct {
	OtherFeesDetails *string `json:"other_fees_details"`
}

type AdditionalInformation struct {
	WebsiteFeeLink string `json:"website_fee_link"`
}

type FinancialData struct {
	Provider              string                `json:"provider"`
	FinancialStatus       FinancialStatus       `json:"financial_status"`
	CreditCertification   CreditCertification   `json:"credit_certification"`
	OtherFees             OtherFees             `json:"other_fees"`
	AdditionalInformation AdditionalInformation `json:"additional_information"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/ConfirmationLetterIssuingServiceFee/CompareProductList"
	payloadTemplate := `{"ProductIdList":"21,49,56,27,52,41,11,62,39,58,53,20,34,30,5,54,43,38,8,64,4,16,44,9,46,7,33,26,60","Page":%d,"Limit":3}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set headers as specified in the curl command
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("cookie", `verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95@AdobeOrg=1; AMCV_F915091E62ED182D0A495F95@AdobeOrg=179643557|MCIDTS|19951|MCMID|53550622918316951353729640026118558196|MCAAMLH-1724305541|3|MCAAMB-1724305541|RKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y|MCOPTOUT-1723707941s|NONE|vVersion|5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.50; visit_time=1763; _ga_NLQFGWVNXN=GS1.1.1724998472.60.1.1725000253.59.0.0`)
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/OtherFee/ConfirmationLetterIssuingServiceFee/CompareProduct")
	req.Header.Set("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", `jGjYP7wTzjQEXOoRLLh1h_dmDwhekANJOyMGhF7ngxhKrJlyOg_uXlaBBE2Hw7ipcSlnYeIyViIEepks0eubAXzYstTb9kQMYysWMhX88301,qj9SnCGTC9a-uq8CMn3rCs88XC4sHFAI0PSx5mi78_N_O5tYMlSoUAJxBrJjgMpDgvdOL737L3b2VDjh223LMAsSV7yGNdRT5LrWJc_yz4A1`)
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	totalPages := DetermineTotalPage(doc)
	if totalPages == 0 {
		log.Fatal("Could not determine the total number of pages")
	}

	var results []FinancialData

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}

		// Set headers as specified in the curl command
		req.Header.Set("accept", "text/plain, */*; q=0.01")
		req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
		req.Header.Set("content-type", "application/json; charset=UTF-8")
		req.Header.Set("cookie", `verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95@AdobeOrg=1; AMCV_F915091E62ED182D0A495F95@AdobeOrg=179643557|MCIDTS|19951|MCMID|53550622918316951353729640026118558196|MCAAMLH-1724305541|3|MCAAMB-1724305541|RKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y|MCOPTOUT-1723707941s|NONE|vVersion|5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.50; visit_time=1763; _ga_NLQFGWVNXN=GS1.1.1724998472.60.1.1725000253.59.0.0`)
		req.Header.Set("origin", "https://app.bot.or.th")
		req.Header.Set("priority", "u=1, i")
		req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/OtherFee/ConfirmationLetterIssuingServiceFee/CompareProduct")
		req.Header.Set("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"macOS"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
		req.Header.Set("verificationtoken", `jGjYP7wTzjQEXOoRLLh1h_dmDwhekANJOyMGhF7ngxhKrJlyOg_uXlaBBE2Hw7ipcSlnYeIyViIEepks0eubAXzYstTb9kQMYysWMhX88301,qj9SnCGTC9a-uq8CMn3rCs88XC4sHFAI0PSx5mi78_N_O5tYMlSoUAJxBrJjgMpDgvdOL737L3b2VDjh223LMAsSV7yGNdRT5LrWJc_yz4A1`)
		req.Header.Set("x-requested-with", "XMLHttpRequest")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making HTTP request: %v", err)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Fatalf("Error parsing HTML: %v", err)
		}
		for i := 1; i <= 3; i++ {
			col := "col" + strconv.Itoa(i)
			var result FinancialData

			// Extract the provider name
			result.Provider = cleanText(doc.Find(fmt.Sprintf("th.%s span", col)).Text())

			// Extract and convert financial status data for this column
			result.FinancialStatus.FinancialCertificationForVisa = extractAndParseCertification(doc.Find(fmt.Sprintf(".attr-VisaEmbassy .%s span", col)))
			result.FinancialStatus.FinancialCertificationForGovernment = extractAndParseCertification(doc.Find(fmt.Sprintf(".attr-GovernantPrivateAgencyInstitution .%s span", col)))
			result.FinancialStatus.FinancialCertificationForAudit = extractAndParseCertification(doc.Find(fmt.Sprintf(".attr-Audit .%s span", col)))

			// Extract credit certification data for this column
			result.CreditCertification.ConditionalCreditCertification = extractAndParseConditionList(doc.Find(fmt.Sprintf(".attr-Conditional .%s span", col)))
			result.CreditCertification.UnconditionalCreditCertification = extractAndParseConditionList(doc.Find(fmt.Sprintf(".attr-Unconditional .%s span", col)))

			// Extract other fees details for this column
			otherFeesText := cleanText(doc.Find(fmt.Sprintf(".attr-other .%s span", col)).Text())
			if otherFeesText != "" {
				result.OtherFees.OtherFeesDetails = &otherFeesText
			}

			// Extract additional information for this column
			result.AdditionalInformation.WebsiteFeeLink = doc.Find(fmt.Sprintf(".attr-Feeurl .%s a.prod-url", col)).AttrOr("href", "")

			results = append(results, result)
		}
		time.Sleep(2 * time.Second)
	}

	// Convert to JSON and print
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	file, err := os.Create("financial_certifi.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}

func extractAndParseCertification(selection *goquery.Selection) Certification {
	text := cleanText(selection.Text())
	value := parseNumericValue(text)
	return Certification{
		OriginalText: text,
		Value:        value,
	}
}

func parseNumericValue(value string) int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(value)
	if match != "" {
		num, err := strconv.Atoi(match)
		if err == nil {
			return num
		}
	}
	return 0 // Return zero if no numeric value is found or conversion fails
}

func extractAndParseConditionList(selection *goquery.Selection) []Certification {
	var certifications []Certification
	text := selection.Text()

	// Remove newlines, extra spaces, and commas
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, ",", "")
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Split text by "-" or numbers to clean up
	splitByDash := strings.Split(text, "-")
	for _, part := range splitByDash {
		part = strings.TrimSpace(part)
		if part != "" {
			certifications = append(certifications, Certification{
				OriginalText: part,
				Value:        parseNumericValue(part),
			})
		}
	}

	// Further split based on numbered list patterns (e.g., "1.", "2.")
	var finalCertifications []Certification
	for _, certification := range certifications {
		splitByNumbers := regexp.MustCompile(`(\d+\.)`).Split(certification.OriginalText, -1)
		for _, part := range splitByNumbers {
			part = strings.TrimSpace(part)
			if part != "" {
				finalCertifications = append(finalCertifications, Certification{
					OriginalText: part,
					Value:        parseNumericValue(part),
				})
			}
		}
	}

	return finalCertifications
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
	text = strings.Join(strings.Fields(text), " ")
	return text
}
