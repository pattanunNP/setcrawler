package utils

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func AddHeader(req *http.Request) {
	// Add headers to the request
	req.Header.Add("accept", "text/plain, */*; q=0.01")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.38; _ctout6672=1; _ga_NLQFGWVNXN=GS1.1.1724821777.47.1.1724821805.32.0.0`)
	req.Header.Add("origin", "https://app.bot.or.th")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://app.bot.or.th/1213/MCPD/FeeApp/CreditFee/CompareProduct")
	req.Header.Add("sec-ch-ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"macOS"`)
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Add("verificationtoken", "dLZqjcdam7cEgePnpg99_6kVdCFqFkyQbCCRj3QV2b1FGm0hWuqRLaRkqfbqoM1H-sjFK9UVNJEFWZVJmP15o3Rx8qmjIdsDHvq1TuWfSlQ1,GuxkO90sOnt2U5WskpIC9Sz9JNDlDDpGPBL4bXdS3WSEyoBtaW3rAK4AAZWEqb0o3GMUnbzoAxlDNqwe0S5ELyu_DXfNBeSr5r4RmG-kysY1")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
}

func CleanText(text string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(text, " "))
}

func SplitText(text, delimeter string) []string {
	parts := strings.Split(text, delimeter)
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
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
