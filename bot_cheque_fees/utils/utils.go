package utils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func AddHeader(req *http.Request) {
	// Set headers
	req.Header.Set("Accept", "text/plain, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", "verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT=\"z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0\"; _uid6672=16B5DEBD.37; _ctout6672=1; visit_time=11; _ga_NLQFGWVNXN=GS1.1.1724811302.45.1.1724811324.38.0.0")
	req.Header.Set("Origin", "https://app.bot.or.th")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://app.bot.or.th/1213/MCPD/FeeApp/ChequeFee/CompareProduct")
	req.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("VerificationToken", "jZuB8Ud2BCv0AAy54UrH0PfUz1cAel3vFRKPTGaOOUYjBjfXr39Vb1ADfdgRUfKZp1wcRIq0SxSvhphifB_TM8Nb9itcqXrJxo_rKh_CRSI1,N3yCdfFo2BB4IDNMKL3RPE44tt0M40SXngT-JFi245TRQQWfejsuZurX0U7PiOREQrywNMQg1pdzx0H0YF0lATg8bFdm43nejrSCwIjyPgo1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

func ExtractFee(doc *goquery.Document, element, col string) string {
	return CleanText(doc.Find("tr." + element + " td." + col).Text())
}

// Extracts an array of fees from the text content based on element and column, split by "-" and newline
func ExtractFeeArray(doc *goquery.Document, element, col string) []string {
	var fees []string
	doc.Find("tr." + element + " td." + col).Each(func(i int, s *goquery.Selection) {
		fee := CleanText(s.Text())
		if fee != "" {
			splitFees := strings.Split(fee, "-")
			for _, f := range splitFees {
				fees = append(fees, CleanText(f))
			}
		}
	})
	return fees
}

// Extracts a link based on element and column
func ExtractLink(doc *goquery.Document, element, col string) string {
	link, _ := doc.Find("tr." + element + " td." + col + " a").Attr("href")
	return link
}

func CleanText(text string) string {
	parts := strings.Fields(text)
	return strings.Join(parts, " ")
}

func SplitByCondition(text string) []string {
	parts := strings.Split(text, "เงื่อนไข")
	var result []string
	for _, part := range parts {
		cleaned := CleanText(part)
		if cleaned != "" {
			result = append(result, "เงื่อนไข"+cleaned)
		}
	}
	return result
}
