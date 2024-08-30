package utils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func AddHeader(req *http.Request) {
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.56; _ctout6672=1; _ga_NLQFGWVNXN=GS1.1.1725014138.63.1.1725014158.40.0.0`)
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/PLoanwithorwithoutCollateralFee/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", "macOS")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", "CnHf4uC4jO6ecbj0JpaCiR98IWnrh7bJMkwXuUAhXfMk_g0c52QzMYBk_cuq0WroG2MIU1ES0DwKQ7WD9sQomBdRQorS_SPDLlMjsIFgltg1,Nl9K4EwzGzfYl0Mi3mN4a-umAOBxD8HjSyT2J3Ht-AAROg1TD8Kp3PjjOXD0Uzat7uL6DAOhuC9dGiBHG7xj3tUcJ1IQYbCoBR8_xX33B-c1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

}

// CleanText removes unwanted characters and trims whitespace
func CleanText(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", "")
	text = strings.Join(strings.Fields(text), " ")
	return text
}

// SplitAndTrim splits a string by a hyphen and trims whitespace from each resulting part
func SplitAndTrim(text string) []string {
	// Split the input text by the hyphen character
	parts := strings.Split(text, "-")

	// Iterate over each part to trim leading and trailing whitespace
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
