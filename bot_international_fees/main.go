package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"internationla_fees/models"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://app.bot.or.th/1213/MCPD/FeeApp/InternationalTransactionFee/CompareProductList"
	payload := []byte(`{"ProductIdList":"122,107,108,97,81,116,13,34,39,37,99,72,20,30,38,118,12,4,52,24,19,85,50,91,104,73,15,109,100,21,66,65,51","Page":1,"Limit":3}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/InternationalTransactionFee/CompareProduct")
	req.Header.Set("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", "rFmM9hbbKpIhSDbHqnSdAEA7EXfxf6BPeusXcaOViG0yZo2bwKfQixWh9515bv11bgOXFI1qVp9pWHu9_lCQjyEF0Bj8dKpqB1ZQ_tdPQrU1,EB_9YBtb5lIDuCmC2AC3vws0nx1WUuVrlfLy0RNs3A0BEw81UrwGMPrKl6NegQP39MfC7HrjOce9dPwWdh719gdhaTjdOjRPCdZyudqINko1")
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	var international_fees []models.Fees

	for i := 1; i <= 3; i++ {
		col := "col" + strconv.Itoa(i)
		provider := doc.Find(fmt.Sprintf("th.%s span", col)).Text()
		provider = CleanText(provider)

		international_fees = append(international_fees, models.Fees{Provider: provider})

	}

	file, err := json.MarshalIndent(international_fees, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile("international_fee.json", file, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Data successfully saved to international_fee.json")
}

func CleanText(text string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(text, " "))
}
