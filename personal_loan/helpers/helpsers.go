package helpers

import (
	"encoding/json"
	"html"
	"net/http"
	"os"
	"personal_loan/pkg/model"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func SaveToJSON(data []model.LoanProduct, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

func CleanString(input string) string {
	trimmed := strings.TrimSpace(strings.ReplaceAll(input, "\n", " "))
	return strings.Join(strings.Fields(trimmed), " ")
}

func GetNullableText(input string) *string {
	trimmed := CleanString(input)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func ParseInt(input string) int {
	cleaned := strings.ReplaceAll(input, ",", "")
	number, _ := strconv.Atoi(regexp.MustCompile(`\d+`).FindString(cleaned))
	return number
}

func ParseLoanAmount(input string) model.LoanAmount {
	// Trim spaces and check if the input is empty or invalid
	input = strings.TrimSpace(input)
	if input == "" || input == "-" {
		// Handle empty or invalid input by returning default values
		return model.LoanAmount{Min: 0, Max: 0}
	}

	// Replace non-numeric characters except for hyphens
	re := regexp.MustCompile(`[^\d\s-]`)
	cleanedInput := re.ReplaceAllString(input, "")

	// Split the cleaned input by "-" to get min and max values
	parts := strings.Split(cleanedInput, "-")

	// Parse the minimum value
	min := ParseInt(parts[0])

	// Parse the maximum value if it exists, otherwise use the min
	var max int
	if len(parts) > 1 {
		max = ParseInt(parts[1])
	} else {
		max = min
	}

	return model.LoanAmount{Min: min, Max: max}
}

func ParseLoanDuration(input string) model.LoanDurationMonths {
	// Trim spaces and handle unexpected format
	input = strings.TrimSpace(input)

	// Check if the input is "-" or empty, and handle it
	if input == "-" || input == "" {
		return model.LoanDurationMonths{MinMonth: 0, MaxMonth: 0}
	}

	// Replace non-numeric characters
	re := regexp.MustCompile(`[^\d\s-]`)
	cleanedInput := re.ReplaceAllString(input, "")

	// Split the cleaned input by "-" to get min and max values
	parts := strings.Split(cleanedInput, "-")

	// Parse the minimum value
	minMonth := ParseInt(parts[0])

	// Parse the maximum value if it exists, otherwise use the minMonth
	var maxMonth int
	if len(parts) > 1 {
		maxMonth = ParseInt(parts[1])
	} else {
		maxMonth = minMonth
	}

	return model.LoanDurationMonths{MinMonth: minMonth, MaxMonth: maxMonth}
}

func FormatCreditLimit(input string) []string {
	input = html.UnescapeString(input)

	input = strings.ReplaceAll(input, "\n", "")
	input = strings.TrimSpace(input)

	parts := strings.Split(input, "- รายได้")

	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			if !strings.HasPrefix(trimmed, "รายได้") {
				trimmed = "รายได้ " + trimmed
			}
			result = append(result, trimmed)
		}
	}
	return result
}

func ParseOptionalField(doc *goquery.Document, selector string) *string {
	text := CleanString(doc.Find(selector).Text())
	if text == "-" || text == "" {
		return nil
	}
	return &text
}

func ParseLinkField(doc *goquery.Document, selector string) *string {
	href, exists := doc.Find(selector).Find("a").Attr("href")
	if exists {
		return &href
	}
	return nil
}

func SetHeaders(req *http.Request) {
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m04nlc8h&sl=0&tt=0"; _uid6672=16B5DEBD.25; _ctout6672=1; visit_time=7; _ga_NLQFGWVNXN=GS1.1.1724410123.29.1.1724412174.45.0.0`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/ProductApp/PersonalLoan/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("Verificationtoken", "iqq1Y2LJoacbaz8ZHPGYLMDGLEjZ3ufpURiA7Dod2r6WisK6Mi-EHz8G5HbGC8ZSmfDfrsQxZkuofF80bKnsIEelVsDpUb3gAoNCpfQmfDI1,CpE7Bvo2sEyVieZ5L7p88hmLZ0xlNf_sAAOLtc8V0q1THAYVZfKncC58hn_6dv4tqV8k2HIez6mzGPUgiqImJCNSGTBN8NDIXJ5nMXxc9rE1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

func SplitAndFormatArray(input string) []string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.TrimSpace(input)

	parts := strings.Split(input, "-")
	var formattedParts []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			formattedParts = append(formattedParts, trimmed)
		}
	}

	return formattedParts
}

func SplitMoneyTransferConditions(input string) []string {
	// Trim and clean up the input
	cleanedInput := strings.TrimSpace(input)

	// Replace all instances of '-' and '/' with a newline separator (you can choose any unique separator)
	cleanedInput = strings.ReplaceAll(cleanedInput, "-", "\n")
	cleanedInput = strings.ReplaceAll(cleanedInput, "/", "\n")

	// Split by newline character to separate conditions
	parts := strings.Split(cleanedInput, "\n")

	// Trim spaces and newline characters for each part and filter out empty strings
	var conditions []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			conditions = append(conditions, trimmed)
		}
	}

	return conditions
}
