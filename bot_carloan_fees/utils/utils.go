package utils

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func AddHeader(req *http.Request) {
	// Set the headers
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.42; visit_time=55503; _ga_NLQFGWVNXN=GS1.1.1724903014.53.0.1724903014.60.0.0`)
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/HirePurchaseFee/CompareProduct")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", "26Ml1XA46AkJbDkhDmoTZot1K37rUT59cf-hdYfNV9adoRoHmk_GO13svlpHj1kD5cO3dWkQC7dXQ95YZec8UP0vU-E4v7cAoHQXf_-wb5Q1,X1IQ7Zo5E0iYY5yZilBT9MGKjmpfCbsqbTUY6ctjYFXflBBfkvpNwOKMK--0iByObZWbuMoFgwpcHtHLQZsunrK9XiZGGtYfSOTF7NLSoPg1")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
}

func CleanText(text string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(text, " "))
}

func SplitTextByDelimiters(text string) *[]string {
	if text == "" {
		return nil
	}

	delimeters := []string{"เงื่อนไข:", "-"}
	for _, delimeter := range delimeters {
		text = strings.ReplaceAll(text, delimeter, "|")
	}

	parts := strings.Split(text, "|")

	cleanedParts := []string{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			cleanedParts = append(cleanedParts, part)
		}
	}
	if len(cleanedParts) == 0 {
		return nil
	}
	return &cleanedParts
}

func SplitText(text string, delimiter string) []string {
	parts := strings.Split(text, delimiter)
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

func SplitTextByNumbers(text string) []string {
	re := regexp.MustCompile(`\d+\.`)
	indices := re.FindAllStringIndex(text, -1)
	if len(indices) == 0 {
		return []string{text}
	}

	var parts []string
	lastIndex := 0
	for _, index := range indices {
		if index[0] > lastIndex {
			part := text[lastIndex:index[0]]
			parts = append(parts, strings.TrimSpace(part))
		}
		lastIndex = index[1]
	}

	parts = append(parts, strings.TrimSpace(text[lastIndex:]))
	return parts
}

func ParseOptionalString(doc *goquery.Document, selector string) *string {
	text := CleanText(doc.Find(selector).Text())
	if text == "" {
		return nil
	}
	return &text
}

func ParseOptionalArray(doc *goquery.Document, selector string) *[]string {
	text := CleanText(doc.Find(selector).Text())
	return SplitTextByDelimiters(text)
}

func DetermineTotalPage(doc *goquery.Document) int {
	totalPages := 1
	doc.Find("#pnlPaging ul.pagination li a").Each(func(i int, s *goquery.Selection) {
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
