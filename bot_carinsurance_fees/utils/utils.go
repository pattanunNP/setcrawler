package utils

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func AddHeader(req *http.Request) {
	// Add headers (Use actual cookie and headers for your request)
	req.Header.Add("accept", "text/plain, */*; q=0.01")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.58; visit_time=1228; _ga_NLQFGWVNXN=GS1.1.1725019349.65.1.1725020596.59.0.0`)
	req.Header.Add("origin", "https://app.bot.or.th")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/TitleLoanFee/CompareProduct")
	req.Header.Add("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"macOS"`)
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Add("verificationtoken", `xuV6ropVxNtfs2a2b1KnXAcqL4b6pz3kQGXSGBpTaSfElop_tC9QrOm9qsmN3QnTukKmzl3wBg0XAM1osRsWCUbSIn6i3Au_5JX_7qvXkOA1,4ob3IH2Vyn3PnDONuCtWdKDjcKamqpmJECNYWHXiR8bgnyKnqwwcTTjh3IgMfWGMYH4zF0U01l-lvmWaFnln4EPKVOhw5FBGFI995ZFMnNw1`)
	req.Header.Add("x-requested-with", "XMLHttpRequest")
}

func CleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", "")
	text = strings.Join(strings.Fields(text), " ")
	return text
}

func SplitText(text string, delimiter string) []string {
	parts := strings.Split(text, delimiter)
	for i, part := range parts {
		parts[i] = CleanText(part)
	}
	return parts
}

func ConvertTextToFloat(text string) float64 {
	text = strings.ReplaceAll(text, ",", "")
	text = strings.TrimSpace(text)
	if text == "" || text == "ไม่มีค่าธรรมเนียม" || text == "ไม่มีบริการ" {
		return 0.0
	}
	var value float64
	fmt.Sscanf(text, "%f", &value)
	return value
}

func ExtractMaxFee(feeTexts []string) float64 {
	var maxFee float64
	for _, feeText := range feeTexts {
		value, _ := strconv.ParseFloat(strings.Fields(feeText)[0], 64)
		if value > maxFee {
			maxFee = value
		}
	}
	return maxFee
}

func ExtractNumbersFromText(text string) []int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(text, -1)
	var numbers []int
	for _, match := range matches {
		if num, err := strconv.Atoi(match); err == nil {
			numbers = append(numbers, num)
		}
	}
	return numbers
}

// ExtractFloatNumbersFromText extracts all floating-point numbers from a text string and returns them as an array of float64.
func ExtractFloatNumbersFromText(text string) []float64 {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	matches := re.FindAllString(text, -1)
	var numbers []float64
	for _, match := range matches {
		if num, err := strconv.ParseFloat(match, 64); err == nil {
			numbers = append(numbers, num)
		}
	}
	return numbers
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
