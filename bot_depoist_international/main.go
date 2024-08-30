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

type BankAccount struct {
	Bank                            string           `json:"Bank"`
	AccountCurrency                 string           `json:"AccountCurrency"`
	MinimumDepositForAccountOpening ValueWithText    `json:"MinimumDepositForAccountOpening"`
	AnnualInterestRate              ValueWithText    `json:"AnnualInterestRate"`
	DepositTerm                     TermWithText     `json:"DepositTerm"`
	MinimumAverageBalance           ValueWithText    `json:"MinimumAverageBalance"`
	FeeIfBalanceBelowMinimum        ValueWithText    `json:"FeeIfBalanceBelowMinimum"`
	FeeIfAccountInactive            CurrencyWithText `json:"FeeIfAccountInactive"`
	Individual                      CurrencyWithText `json:"Individual"`
	Corporate                       CurrencyWithText `json:"Corporate"`
}

type ValueWithText struct {
	OriginalText string  `json:"OriginalText"`
	Value        *string `json:"Value,omitempty"`
}

type TermWithText struct {
	OriginalText string `json:"OriginalText"`
	Months       *int   `json:"Months,omitempty"`
}

type CurrencyWithText struct {
	OriginalText string  `json:"OriginalText"`
	Value        *string `json:"Value,omitempty"`
	Currency     *string `json:"Currency,omitempty"`
}

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FCDInterestAndFeeRateApp/Search/SearchProductInformation"
	payloadTemplate := `{"FICodeList":"","ORG_IP_ID_List":null,"AR_TP_ID_List":null,"INT_RATE_TYPE_List":null,"INACT_PERIOD_List":null,"INACT_FEE":null,"CCY_ID_List":null,"Page":%d,"DisplayOrder":"1"}`

	initialPayload := fmt.Sprintf(payloadTemplate, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(initialPayload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	setHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Non-OK HTTP status: %s", resp.Status)
	}

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
	// totalPages := 3

	var bankAccounts []BankAccount

	for page := 1; page <= totalPages; page++ {
		fmt.Printf("Processing page: %d/%d\n", page, totalPages)
		payload := fmt.Sprintf(payloadTemplate, page)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		if err != nil {
			fmt.Printf("Error creating request for page %d: %v\n", page, err)
			continue
		}

		setHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request for page %d: %v\n", page, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Non-OK HTTP status for page %d: %s\n", page, resp.Status)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body for page %d: %v\n", page, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			fmt.Printf("Error parsing HTML for page %d: %v\n", page, err)
			continue
		}

		doc.Find("tr").Each(func(index int, item *goquery.Selection) {
			bank := strings.TrimSpace(item.Find(".prod-bank").Text())
			accountCurrency := strings.TrimSpace(item.Find("td").Eq(2).Text())

			// Validate essential fields to avoid including empty entries
			if bank == "" || accountCurrency == "" {
				return // Skip empty or invalid entries
			}

			minDepositText := strings.TrimSpace(item.Find("td").Eq(3).Contents().First().Text())
			minDepositValue := removeCommas(minDepositText)

			annualInterestRateText := strings.TrimSpace(item.Find("td").Eq(5).Text())
			annualInterestRateValue := removeCommas(annualInterestRateText)

			depositTermText := strings.TrimSpace(item.Find("td").Eq(6).Text())
			months := extractMonths(depositTermText)

			minAvgBalanceText := strings.TrimSpace(item.Find("td").Eq(7).Text())
			minAvgBalanceValue := removeCommas(minAvgBalanceText)

			feeBelowMinText := strings.TrimSpace(item.Find("td").Eq(8).Text())
			feeBelowMinValue := removeCommas(feeBelowMinText)

			feeInactiveText := strings.TrimSpace(item.Find("td").Eq(9).Text())
			feeInactiveValue, feeInactiveCurrency := splitCurrency(feeInactiveText)

			individualText := strings.TrimSpace(item.Find("td").Eq(10).Text())
			individualValue, individualCurrency := splitCurrency(individualText)

			corporateText := "" // Assuming you do not have a corporate column
			corporateValue := ""

			account := BankAccount{
				Bank:            bank,
				AccountCurrency: accountCurrency,
				MinimumDepositForAccountOpening: ValueWithText{
					OriginalText: minDepositText,
					Value:        &minDepositValue,
				},
				AnnualInterestRate: ValueWithText{
					OriginalText: annualInterestRateText,
					Value:        &annualInterestRateValue,
				},
				DepositTerm: TermWithText{
					OriginalText: depositTermText,
					Months:       months,
				},
				MinimumAverageBalance: ValueWithText{
					OriginalText: minAvgBalanceText,
					Value:        &minAvgBalanceValue,
				},
				FeeIfBalanceBelowMinimum: ValueWithText{
					OriginalText: feeBelowMinText,
					Value:        &feeBelowMinValue,
				},
				FeeIfAccountInactive: CurrencyWithText{
					OriginalText: feeInactiveText,
					Value:        &feeInactiveValue,
					Currency:     &feeInactiveCurrency,
				},
				Individual: CurrencyWithText{
					OriginalText: individualText,
					Value:        &individualValue,
					Currency:     &individualCurrency,
				},
				Corporate: CurrencyWithText{
					OriginalText: corporateText,
					Value:        &corporateValue,
					Currency:     nil,
				},
			}

			bankAccounts = append(bankAccounts, account)
		})
	}

	jsonData, err := json.MarshalIndent(bankAccounts, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	err = os.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("Data saved to output.json")
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

func setHeaders(req *http.Request) {
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("cookie", `verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.60; visit_time=1472; _ga_NLQFGWVNXN=GS1.1.1725024405.66.1.1725026441.60.0.0`)
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/FCDInterestAndFeeRateApp")
	req.Header.Set("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", "4r4xNjC644P-I_hp2Bum2TtvVh0gVg5eY5Kao7miivJdRHgNRB_quTKkTLmBQsJl3CqAedlpa0KR5DswlHp1BG7ukFSH4nLGmq1sVx1HaQY1,TnMfKf12vw6o2YxnF8GpUp9_8o__3KfDlTXMUqV8pMZna391EzwV0eWQuOydM_1wbmM7anQuD5G3sAVTqCPTHEW83pb-fxeXnMbGkG8Q6nM1")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
}

func removeCommas(value string) string {
	return strings.ReplaceAll(value, ",", "")
}

func extractMonths(term string) *int {
	term = strings.TrimSpace(term)
	if strings.Contains(term, "เดือน") {
		term = strings.Replace(term, "เดือน", "", 1)
		term = strings.TrimSpace(term)
		if months, err := strconv.Atoi(term); err == nil {
			return &months
		}
	}
	return nil
}

func splitCurrency(value string) (string, string) {
	parts := strings.Fields(value)
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return value, ""
}
