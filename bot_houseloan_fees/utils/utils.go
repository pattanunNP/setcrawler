package utils

import (
	"houseLoan_fees/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func AddHeader(req *http.Request) {
	// Setting the necessary headers
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.45; visit_time=7364; _ga_NLQFGWVNXN=GS1.1.1724931231.56.1.1724931232.59.0.0`)
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/homeloanFee/CompareProduct")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", `yPWJNwR_QEP9E8zXYUWUjBzIpnYkzmr1XRjnxBpLPj2l0D0XRxtDmb5T0ZVqaSkaYlM_X8p47Tqpm1df2wi_bE4XpJ4nfVJgCh3euZSKhcA1,bbqurLXSvvkwH-8BgEXKbLAn0usgNAutyLXwkZbkwNyOXys_OUhk3a7Wxwb5oIbbXF04c_-v0fMMQeO7betnL2-RcdPAzcITcTcKGLv0Ksw1`)
	req.Header.Set("x-requested-with", "XMLHttpRequest")
}

func CleanText(text string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(text, " "))
}

func CleanTextArray(text string) []string {
	text = strings.ReplaceAll(text, ",", "")

	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, " ", "")

	parts := strings.Split(text, "-")

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func SplitAndCleanText(text, delimeter string) []string {
	parts := strings.Split(text, delimeter)
	var cleanedParts []string
	for _, part := range parts {
		cleaned := CleanTextArray(part)
		cleanedParts = append(cleanedParts, cleaned...)
	}
	return cleanedParts
}

func NullEmpty(text string) *string {
	cleanedText := strings.Join(CleanTextArray(text), " ")
	if cleanedText == "" {
		return nil
	}
	return &cleanedText
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

func ParseFeeDetail(text string) models.FeeDetail {
	cleanedText := CleanTextArray(text)
	parseValue := parseNumericValue(strings.Join(cleanedText, " "))
	minFee, maxFee := extractMinMaxFee(strings.Join(cleanedText, " "))
	return models.FeeDetail{
		Original: cleanedText,
		Numeric:  parseValue,
		MinFee:   minFee,
		MaxFee:   maxFee,
	}
}

func extractMinMaxFee(text string) (minFee, maxFee float64) {
	minFee, maxFee = 0, 0
	reMin := regexp.MustCompile(`ไม่ต่ำกว่าฉบับละ (\d+) บาท`)
	reMax := regexp.MustCompile(`สูงสุดไม่เกิน (\d+) บาท`)
	minMatch := reMin.FindStringSubmatch(text)
	maxMatch := reMax.FindStringSubmatch(text)

	if len(minMatch) > 1 {
		minValue, err := strconv.ParseFloat(minMatch[1], 64)
		if err == nil {
			minFee = minValue
		}
	}

	if len(maxMatch) > 1 {
		maxValue, err := strconv.ParseFloat(maxMatch[1], 64)
		if err == nil {
			maxFee = maxValue
		}
	}

	return minFee, maxFee
}

func parseNumericValue(text string) *float64 {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	match := re.FindString(text)
	if match != "" {
		value, err := strconv.ParseFloat(match, 64)
		if err == nil {
			return &value
		}
	}
	return nil
}
