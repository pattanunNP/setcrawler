package utils

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func AddHeader(req *http.Request) {
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.46; _ctout6672=1; visit_time=728; _ga_NLQFGWVNXN=GS1.1.1724931231.56.1.1724932645.59.0.0`)
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/InternationalTransactionFee/CompareProduct")
	req.Header.Set("sec-ch-ua", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", `_I3hbpPk_RQpCCswztEuEVgyf07AMW7EyyLyfxDYY7pQjW0mW5TzcbmeHK6ZieefsNcqRa4B1avjVi561hD8ZZl72zqxF353xY5KJOwdMWQ1,i4XoNeFHUHioilxc130gSxUlNufrnlSzDiVyHtYzusl-QAYmUQdyapZhBKuZPXWgwAMEz68oLQUHGTMx8RImRW78SLDCH81W74TtsMPv-WM1`)
	req.Header.Set("x-requested-with", "XMLHttpRequest")
}

func CleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	text = strings.Join(strings.Fields(text), " ")
	return text
}

func SplitAndCleanText(text, delimeter string) []string {
	parts := strings.Split(text, delimeter)
	var cleanedParts []string
	for _, part := range parts {
		cleaned := CleanText(part)
		if cleaned != "" {
			cleanedParts = append(cleanedParts, cleaned)
		}
	}
	return cleanedParts
}

func CleanConditionText(text string) string {
	// Replace newlines and tabs with a space
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	// Remove commas
	text = strings.ReplaceAll(text, ",", "")

	// Trim spaces
	text = strings.TrimSpace(text)

	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	return text
}

func ProcessTextWithPattern(text string) []string {
	text = CleanConditionText(text)

	re := regexp.MustCompile(`(\d+\.)`)
	parts := re.Split(text, -1)

	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func ExtractTextInsideElement(doc *goquery.Document, selector string) string {
	var textParts []string
	doc.Find(selector).Each(func(i int, td *goquery.Selection) {
		td.Find("span").Each(func(j int, span *goquery.Selection) {
			text := strings.TrimSpace(span.Text())
			if text != "" {
				textParts = append(textParts, text)
			}
		})
	})
	return strings.Join(textParts, " ")
}

// Function to extract compensation fee details from a <td> element
func ExtractCompensationFeeFromTd(doc *goquery.Document, selector string) string {
	var compensationFee string
	doc.Find(selector).Each(func(i int, td *goquery.Selection) {
		td.Find("span").Each(func(j int, span *goquery.Selection) {
			if strings.Contains(span.Text(), "ค่าธรรมเนียมชดเชยอัตราแลกเปลี่ยน") {
				compensationFee = strings.TrimSpace(span.Text())
			}
		})
	})
	return compensationFee
}

func ExtractFee(text string) string {
	// Split the text by the known delimiter
	parts := strings.Split(text, "ค่าธรรมเนียมชดเชยอัตราแลกเปลี่ยน")
	return strings.TrimSpace(parts[0])
}

func ExtractCompensationFee(text string) string {
	// Look for compensation fee details
	if strings.Contains(text, "ค่าธรรมเนียมชดเชยอัตราแลกเปลี่ยน") {
		return "ไม่มีค่าธรรมเนียม"
	}
	return "ไม่มีบริการ"
}

func ExtractFeeType(text string) string {
	if strings.Contains(text, "ตามอัตราที่กำหนด") {
		return "ตามอัตราที่กำหนด"
	}
	return "ไม่มีบริการ"
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
